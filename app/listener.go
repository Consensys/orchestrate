package app

import (
	"context"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"github.com/ethereum/go-ethereum/common/hexutil"	
	"github.com/spf13/viper"
	"github.com/spf13/pflag"
	"gitlab.com/ConsenSys/client/fr/core-stack/core.git"
	tracepb "gitlab.com/ConsenSys/client/fr/core-stack/core.git/protobuf/trace"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/tx-listener"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-listener.git/app/infra"
	"gitlab.com/ConsenSys/client/fr/core-stack/worker/tx-listener.git/app/worker"
)

// ListenerHandler handler that listen to tx-listener messages
type ListenerHandler struct {
	app *App

	startPositions  map[string]*StartingPosition
	defaultPosition *StartingPosition

	cleanOnce *sync.Once
	worker    *core.Worker
	logger    *log.Entry

	cfg listener.Config

	listener listener.TxListener
}

// Setup creates listener
func (l *ListenerHandler) Setup() {
	l.worker = worker.CreateWorker(l.app.infra)

	// Pipe sarama message channel into worker
	in := make(chan interface{})
	go func() {
		// Pipe channels for interface compatibility
		pipeLoop:
		for {
			select {
			case <-l.app.ctx.Done():
				break pipeLoop

			case i, ok := <-l.listener.Receipts():
				if !ok {
					break pipeLoop
				} else {
					in <- i
				}
			}
		}
		close(in)
	}()

	// Run worker
	go l.worker.Run(in)
	
	// Start draining errors
	go func() {
		for err := range l.listener.Errors() {
			log.WithFields(log.Fields{
				"Chain": err.ChainID.Text(16),
			}).WithError(err).Errorf("tx-listener: got error")
		}
	}()
	
	// Start draining blocks
	go func() {
		for block := range l.listener.Blocks() {
			log.WithFields(log.Fields{
				"BlockHash":   block.Header().Hash().Hex(),
				"BlockNumber": block.Header().Number,
				"Chain":       block.ChainID.Text(16),
			}).Debugf("tx-listener: got new block")
		}
	}()

	log.Infof("tx-listener: ready to listen")
}

// Listen start listening
func (l *ListenerHandler) Listen() {
	for _, chainID := range l.app.infra.Mec.Networks(context.Background()) {
		// Start listening
		position, err := l.getStartingPosition(chainID)
		if err != nil {

		}
		l.listener.Listen(chainID, position.BlockNumber, position.TxIndex, l.cfg)
	}
	// Wait for worker to be done
	<-l.worker.Done()

	// Close listener
	l.listener.Close()
}

func (l *ListenerHandler) getStartingPosition(chainID *big.Int) (*StartingPosition, error) {
	position, ok := l.startPositions[hexutil.EncodeBig(chainID)]
	if !ok {
		position = l.defaultPosition
	}

	if position.BlockNumber != -2 {
		return position, nil
	}

	// BlockNumber == -2 means we should start listening from position of last produce message
	// Compute output topic
	outTopic := fmt.Sprintf("%v-%v", viper.GetString("worker.out"), chainID.Text(16))

	// Retrieve last record
	lastRecord, err := infra.GetLastRecord(l.app.infra.SaramaClient, outTopic, 0)
	if err != nil {
		return nil, err
	}

	if lastRecord == nil {
		// If we have never produced then we start from genesis
		return &StartingPosition{}, nil
	}
	
	// Parse last record using protobuffer
	var pb tracepb.Trace
	err = proto.Unmarshal(lastRecord.Value, &pb)
	if err != nil {
		return nil, err
	}
	
	return &StartingPosition{int64(pb.Receipt.BlockNumber), int64(pb.Receipt.TxIndex)+1}, nil
}

// TranslateBlockNumber translate a starting block number into its integer value
func TranslateBlockNumber(blockNumber string) (int64, error) {
	switch blockNumber {
	case "genesis":
		return 0, nil
	case "latest":
		return -1, nil
	case "oldest":
		return -2, nil
	default:
		res, err := strconv.ParseInt(blockNumber, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("%q is an invalid starting blockNumber expected 'latest', 'oldest', 'genesis' or an integer", blockNumber)
		}
		return res, nil
	}
}

var (
	positionRegexp  = `(?P<chain>0x[a-fA-F0-9]+):(?P<blockNumber>genesis|latest|oldest|\d+)(-(?P<txIndex>\d+))?`
	positionPattern = regexp.MustCompile(positionRegexp)
)

// StartingPosition is an helpful type for storing a starting position
type StartingPosition struct {
	BlockNumber int64
	TxIndex     int64
}

// ParseStartingPosition extract chainID, blockNumber and TxIndex from a formatted starting position string
func ParseStartingPosition(position string) (string, *StartingPosition, error) {
	match := positionPattern.FindStringSubmatch(position)
	if len(match) != 5 {
		return "", nil, fmt.Errorf("Could not parse position (expected format %q): %v ", position, positionRegexp)
	}

	blockNumber, err := TranslateBlockNumber(match[2])
	if err != nil {
		return "", nil, fmt.Errorf("Could not parse position (expected format %q): %v ", position, positionRegexp)
	}

	if match[4] == "" {
		return match[1], &StartingPosition{blockNumber, 0}, nil
	}

	txIndex, err := strconv.ParseInt(match[4], 10, 64)
	if err != nil {
		return "", nil, fmt.Errorf("Could not parse position (expected format %q): %v ", position, positionRegexp)
	}
	return match[1], &StartingPosition{blockNumber, txIndex}, nil
}

// ParseStartingPositions parse starting positions
func ParseStartingPositions(positions []string) (map[string]*StartingPosition, error) {
	m := make(map[string]*StartingPosition)
	for _, position := range positions {
		chain, position, err := ParseStartingPosition(position)
		if err != nil {
			return nil, err
		}
		m[chain] = position
	}
	return m, nil
}

func initListener(app *App) {
	positions, err := ParseStartingPositions(viper.GetStringSlice("listener.start"))
	if err != nil {
		log.WithError(err).Fatalf("tx-listener: could not parse starting positions")
	}

	defaultPosition, err := TranslateBlockNumber(viper.GetString("listener.start.default"))
	if err != nil {
		log.WithError(err).Fatalf("tx-listener: could not parse default starting position")
	}

	config := listener.NewConfig()
	config.BlockCursor.Backoff = viper.GetDuration("listener.block.backoff")
	config.BlockCursor.Limit = uint64(viper.GetInt64("listener.block.limit"))
	config.TxListener.Return.Blocks = true
	config.TxListener.Return.Errors = true

	app.listener = &ListenerHandler{
		app:             app,
		startPositions:  positions,
		defaultPosition: &StartingPosition{defaultPosition,0},
		cleanOnce:       &sync.Once{},
		cfg: config,
		listener: listener.NewTxListener(listener.NewEthClient(app.infra.Mec, config)),
	}
	app.listener.Setup()
}

// InitFlags register flags for listener
func InitFlags(f *pflag.FlagSet) {
	ListenerBlockBackoff(f)
	ListenerBlockLimit(f)
	ListenerTrackerDepth(f)
	ListenerStartDefault(f)
	ListenerStart(f)
}

var (
	listenerBlockBackoffFlag     = "listener-block-backoff"
	listenerBlockBackoffViperKey = "listener.block.backoff"
	listenerBlockBackoffDefault  = time.Second
	listenerBlockBackoffEnv  = "LISTENER_BLOCK_BACKOFF"
)

// ListenerBlockBackoff register flag for Listener Block backoff
func ListenerBlockBackoff(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Backoff time to wait before retrying after failing to find a mined block
Environment variable: %q`, listenerBlockBackoffEnv)
	f.Duration(listenerBlockBackoffFlag, listenerBlockBackoffDefault, desc)
	viper.BindPFlag(listenerBlockBackoffViperKey, f.Lookup(listenerBlockBackoffFlag))
	viper.BindEnv(listenerBlockBackoffViperKey, listenerBlockBackoffEnv)
}

var (
	listenerBlockLimitFlag     = "listener-block-limit"
	listenerBlockLimitViperKey = "listener.block.limit"
	listenerBlockLimitDefault  = int64(40)
	listenerBlockLimitEnv  = "LISTENER_BLOCK_LIMIT"
)

// ListenerBlockLimit register flag for Listener Block limit
func ListenerBlockLimit(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Limit number of block that can be prefetched while listening
Environment variable: %q`, listenerBlockLimitEnv)
	f.Int64(listenerBlockLimitFlag, listenerBlockLimitDefault, desc)
	viper.BindPFlag(listenerBlockLimitViperKey, f.Lookup(listenerBlockLimitFlag))
	viper.BindEnv(listenerBlockLimitViperKey, listenerBlockLimitEnv)
}

var (
	listenerTrackerDepthFlag     = "listener-tracker-depth"
	listenerTrackerDepthViperKey = "listener.tracker.depth"
	listenerTrackerDepthDefault  = int64(5)
	listenerTrackerDepthEnv  = "LISTENER_TRACKER_DEPTH"
)

// ListenerTrackerDepth register flag for Listener Tracker Depth
func ListenerTrackerDepth(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Depth at which we consider a block final (to avoid falling into a re-org)
Environment variable: %q`, listenerTrackerDepthEnv)
	f.Int64(listenerTrackerDepthFlag, listenerTrackerDepthDefault, desc)
	viper.BindPFlag(listenerTrackerDepthViperKey, f.Lookup(listenerTrackerDepthFlag))
	viper.BindEnv(listenerTrackerDepthViperKey, listenerTrackerDepthEnv)
}

var (
	listenerStartDefaultFlag     = "listener-start-default"
	listenerStartDefaultViperKey = "listener.start.default"
	listenerStartDefaultDefault  = "oldest"
	listenerStartDefaultEnv  = "LISTENER_START_DEFAULT"
)

// ListenerStartDefault register flag for Listener Start Default
func ListenerStartDefault(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Default block position listener should start listening from (one of 'latest', 'oldest', 'genesis')
Environment variable: %q`, listenerStartDefaultEnv)
	f.String(listenerStartDefaultFlag, listenerStartDefaultDefault, desc)
	viper.BindPFlag(listenerStartDefaultViperKey, f.Lookup(listenerStartDefaultFlag))
	viper.BindEnv(listenerStartDefaultViperKey, listenerStartDefaultEnv)
}

var (
	listenerStartFlag     = "listener-start"
	listenerStartViperKey = "listener.start"
	listenerStartDefault  = []string{}
	listenerStartEnv  = "LISTENER_START"
)

// ListenerStart register flag for Listener Start Position
func ListenerStart(f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Position listener should start listening from (format <chainID>:<blockNumber>-<txIndex> or <chainID>:<blockNumber>) (e.g. 0x2a:2348721-5 or 0x3:latest)
Environment variable: %q`, listenerStartEnv)
	f.StringSlice(listenerStartFlag, listenerStartDefault, desc)
	viper.BindPFlag(listenerStartViperKey, f.Lookup(listenerStartFlag))
	viper.BindEnv(listenerStartViperKey, listenerStartEnv)
}
