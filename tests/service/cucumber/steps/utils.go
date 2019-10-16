package steps

import (
	"fmt"
)

func LongKeyOf(topic, scenarioID, metadataID string) string {
	return fmt.Sprintf(
		"%v/%v/%v",
		scenarioID,
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
