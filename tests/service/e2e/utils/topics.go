package utils

import (
	"github.com/ConsenSys/orchestrate/pkg/broker/sarama"
	"github.com/ConsenSys/orchestrate/tests/utils"
)

var TOPICS = map[string]string{
	utils.TxSenderTopicKey:  sarama.TxSenderViperKey,
	utils.TxDecodedTopicKey: sarama.TxDecodedViperKey,
	utils.TxRecoverTopicKey: sarama.TxRecoverViperKey,
}
