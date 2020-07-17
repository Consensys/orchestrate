package scheduler

import (
	"sync"

	txscheduler "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/client"

	log "github.com/sirupsen/logrus"
)

const component = "faucet.sarama"

var (
	fct      *Faucet
	initOnce = &sync.Once{}
)

// Init initializes Faucet
func Init() {
	initOnce.Do(func() {
		if fct != nil {
			return
		}

		txscheduler.Init()
		fct = NewFaucet(txscheduler.GlobalClient())

		log.Info("faucet: ready")
	})
}

// GlobalFaucet returns global Sarama Faucet
func GlobalFaucet() *Faucet {
	return fct
}
