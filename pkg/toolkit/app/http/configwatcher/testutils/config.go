package testutils

import (
	"fmt"
	"sort"
)

type Message struct {
	Name string
	Conf string
}

func (msg *Message) ProviderName() string {
	return msg.Name
}

func (msg *Message) Configuration() interface{} {
	return msg.Conf
}

func MergeConfiguration(confs map[string]interface{}) interface{} {
	var merged []string
	for k, v := range confs {
		merged = append(merged, fmt.Sprintf("%v@%v", v.(string), k))
	}
	sort.Strings(merged)
	return merged
}
