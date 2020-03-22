package dynamic

import (
	traefikdynamic "github.com/containous/traefik/v2/pkg/config/dynamic"
)

// Message holds configuration information exchanged between parts of traefik.
type Message struct {
	providerName  string
	configuration *Configuration
}

func NewMessage(providerName string, conf *Configuration) *Message {
	return &Message{
		providerName:  providerName,
		configuration: conf,
	}
}

func (msg *Message) ProviderName() string {
	return msg.providerName
}

func (msg *Message) Configuration() interface{} {
	return msg.configuration
}

func FromTraefikMessage(traefikMsg *traefikdynamic.Message) *Message {
	return &Message{
		providerName:  traefikMsg.ProviderName,
		configuration: FromTraefikConfig(traefikMsg.Configuration),
	}
}

func ToTraefikMessage(msg *Message) *traefikdynamic.Message {
	return &traefikdynamic.Message{
		ProviderName:  msg.providerName,
		Configuration: ToTraefikConfig(msg.configuration),
	}
}
