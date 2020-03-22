package static

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

// +k8s:deepcopy-gen=true

type Interceptor struct {
	Logrus     *Logrus     `json:"logrus,omitempty" toml:"logrus,omitempty" yaml:"logrus,omitempty"`
	Prometheus *Prometheus `json:"prometheus,omitempty" toml:"prometheus,omitempty" yaml:"prometheus,omitempty"`
	Recovery   *Recovery   `json:"recovery,omitempty" toml:"recovery,omitempty" yaml:"recovery,omitempty"`
	Auth       *Auth       `json:"auth,omitempty" toml:"auth,omitempty" yaml:"auth,omitempty"`
	Tracing    *Tracing    `json:"tracing,omitempty" toml:"tracing,omitempty" yaml:"tracing,omitempty"`
	Tags       *Tags       `json:"tags,omitempty" toml:"tags,omitempty" yaml:"tags,omitempty"`
	Error      *Error      `json:"error,omitempty" toml:"error,omitempty" yaml:"error,omitempty"`
}

func (i *Interceptor) Field() (interface{}, error) {
	return utils.ExtractField(i)
}

// +k8s:deepcopy-gen=true

type Logrus struct {
	Fields map[string]string `json:"fields,omitempty" toml:"fields,omitempty" yaml:"fields,omitempty"`
}

// +k8s:deepcopy-gen=true

type Prometheus struct{}

// +k8s:deepcopy-gen=true

type Recovery struct{}

// +k8s:deepcopy-gen=true

type Auth struct{}

// +k8s:deepcopy-gen=true

type Tracing struct {
	TraceHeaderName string `json:"traceHeaderName,omitempty" toml:"traceHeaderName,omitempty" yaml:"traceHeaderName,omitempty"`
}

// +k8s:deepcopy-gen=true

type Tags struct{}

// +k8s:deepcopy-gen=true

type Error struct{}
