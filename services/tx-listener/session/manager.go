package session

import (
	"context"
	"fmt"
	"sync"

	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	ethclientutils "github.com/consensys/orchestrate/pkg/toolkit/ethclient/utils"
	"github.com/consensys/orchestrate/pkg/utils"
	"github.com/consensys/orchestrate/services/tx-listener/dynamic"
	provider "github.com/consensys/orchestrate/services/tx-listener/providers"
)

const component = "tx-listener.session-manager"

type cancelableSession struct {
	session Session
	cancel  func()
}

type Manager struct {
	// Protect sessions mapping
	mux      *sync.RWMutex
	sessions map[string]*cancelableSession

	// Wait Group to keep track of running sessions
	wg *sync.WaitGroup

	// Envelope used to create listening sessions
	builder Builder

	// Configuration
	currentConfiguration *dynamic.Configuration
	msgInput             chan *dynamic.Message

	// Configuration provider
	provider provider.Provider

	commands chan *Command

	errors chan error

	logger *log.Logger
}

func NewManager(sessionBuilder Builder, prvdr provider.Provider) *Manager {
	return &Manager{
		mux:                  &sync.RWMutex{},
		sessions:             make(map[string]*cancelableSession),
		wg:                   &sync.WaitGroup{},
		currentConfiguration: &dynamic.Configuration{},
		msgInput:             make(chan *dynamic.Message),
		builder:              sessionBuilder,
		provider:             prvdr,
		commands:             make(chan *Command),
		errors:               make(chan error),
		logger:               log.NewLogger().SetComponent(component),
	}
}

func (m *Manager) Errors() <-chan error {
	return m.errors
}

func (m *Manager) Run(ctx context.Context) error {
	m.run(ctx)
	return nil
}

func (m *Manager) run(ctx context.Context) {
	defer func() {
		m.wg.Wait()
		m.logger.Info("service stopped")
	}()

	utils.InParallel(
		// Start provider and close input channel when Provider is done
		func() {
			m.listenProvider(ctx)
		},
		// Listen configuration
		func() { m.listenConfiguration() },
		// Listen commands
		func() { m.listenCommands(ctx) },
		func() {
			<-ctx.Done()
			m.logger.WithError(ctx.Err()).Info("service finished gracefully")
		},
	)
}

func (m *Manager) listenProvider(ctx context.Context) {
	m.logger.WithField("provider", fmt.Sprintf("%T", m.provider)).Debug("Starting provider")
	err := m.provider.Run(ctx, m.msgInput)
	if err != nil {
		m.logger.WithError(err).Error("error while listening provider")
	}
	close(m.msgInput)
}

func (m *Manager) listenConfiguration() {
	for msg := range m.msgInput {
		commands := CompareConfiguration(m.currentConfiguration, msg.Configuration)
		if len(commands) > 0 {
			for _, command := range commands {
				m.commands <- command
			}
			m.currentConfiguration = msg.Configuration
		}
	}
	close(m.commands)
}

func (m *Manager) listenCommands(ctx context.Context) {
	for command := range m.commands {
		m.executeCommand(ctx, command)
	}
}

func (m *Manager) executeCommand(ctx context.Context, command *Command) {
	ctx = log.WithFields(ctx,
		log.Field("chain", command.Chain.UUID),
		log.Field("tenant_id", command.Chain.TenantID),
		log.Field("owner_id", command.Chain.OwnerID))
	switch command.Type {
	case START:
		m.runSession(ethclientutils.RetryConnectionError(ctx, true), command.Chain)
	case STOP:
		m.stopSession(ctx, command.Chain)
	case UPDATE:
		m.stopSession(ctx, command.Chain)
		m.runSession(ctx, command.Chain)
	default:
		m.logger.WithContext(ctx).WithField("cmd_type", command.Type).Errorf("unknown command")
	}
}

func (m *Manager) runSession(ctx context.Context, chain *dynamic.Chain) {
	logger := m.logger.WithContext(ctx)
	// Build session
	s, err := m.builder.NewSession(chain)
	if err != nil {
		logger.WithError(err).Errorf("failed to create a new session")
		return
	}

	// Make session cancelable session so we can stop it later on
	ctx, cancel := context.WithCancel(ctx)
	sess := &cancelableSession{
		session: s,
		cancel:  cancel,
	}

	// Add session
	m.addSession(chain.UUID, sess)

	// Start goroutine to run session
	m.wg.Add(1)
	go func() {
		logger.Info("listener session started")
		err := sess.session.Run(ctx)
		m.removeSession(chain.UUID)
		if err != nil && ctx.Err() == nil {
			logger.WithError(err).Error("failed to remove session")
		}
		m.logger.WithField("chain", chain.UUID).Info("session stopped")
		m.wg.Done()
	}()
}

func (m *Manager) stopSession(ctx context.Context, chain *dynamic.Chain) {
	logger := m.logger.WithContext(ctx)
	if sess, ok := m.getSession(chain.UUID); ok {
		logger.WithField("chain", chain.UUID).Debug("stopping session")
		sess.cancel()
		return
	}

	logger.WithField("chain", chain.UUID).Warn("trying to stop a not exiting session")
}

func (m *Manager) addSession(key string, sess *cancelableSession) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.sessions[key] = sess
}

func (m *Manager) removeSession(key string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	delete(m.sessions, key)
}

func (m *Manager) getSession(key string) (sess *cancelableSession, ok bool) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	sess, ok = m.sessions[key]
	return
}
