package utils

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/broker/sarama"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/tests/utils"
)

var TOPICS = map[string]string{
	utils.TxSenderTopicKey:  sarama.TxSenderViperKey,
	utils.TxDecodedTopicKey: sarama.TxDecodedViperKey,
	utils.TxRecoverTopicKey: sarama.TxRecoverViperKey,
}
