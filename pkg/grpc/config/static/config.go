package static

import "time"

// +k8s:deepcopy-gen=true

type Configuration struct {
	Options      *Options       `json:"options,omitempty" toml:"options,omitempty" yaml:"options,omitempty"`
	Interceptors []*Interceptor `json:"interceptors,omitempty" toml:"interceptors,omitempty" yaml:"interceptors,omitempty"`
	Services     *Services      `json:"services,omitempty" toml:"services,omitempty" yaml:"services,omitempty"`
}

type Options struct {
	ConnectionTimeout    time.Duration `json:"connectionTimeout,omitempty" toml:"connectionTimeout,omitempty" yaml:"connectionTimeout,omitempty"`
	HeaderTableSize      uint32        `json:"headerTableSize,omitempty" toml:"headerTableSize,omitempty" yaml:"headerTableSize,omitempty"`
	MaxConcurrentStreams uint32        `json:"maxConcurrentStreams,omitempty" toml:"maxConcurrentStreams,omitempty" yaml:"maxConcurrentStreams,omitempty"`
	MaxHeaderListSize    uint32        `json:"maxHeaderListSize,omitempty" toml:"maxHeaderListSize,omitempty" yaml:"maxHeaderListSize,omitempty"`
	MaxRecvMsgSize       int           `json:"maxRecvMsgSize,omitempty" toml:"maxRecvMsgSize,omitempty" yaml:"maxRecvMsgSize,omitempty"`
	MaxSendMsgSize       int           `json:"maxSendMsgSize,omitempty" toml:"maxSendMsgSize,omitempty" yaml:"maxSendMsgSize,omitempty"`
}
