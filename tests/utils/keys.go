package utils

import (
	"fmt"
)

const TxDecodedTopicKey = "tx.decoded"
const TxRecoverTopicKey = "tx.recover"
const TxSenderTopicKey = "tx.sender"

func LongKeyOf(topic, metadataID string) string {
	return fmt.Sprintf(
		"%v/%v",
		topic,
		metadataID,
	)
}

func ShortKeyOf(topic, scenarioID string) string {
	return fmt.Sprintf(
		"%v/%v",
		scenarioID,
		topic,
	)
}
