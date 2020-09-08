package utils

import (
	"fmt"
)

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
