package utils

type GlobalData struct {
	Nodes Nodes `json:"nodes"`
}

type Nodes struct {
	BesuOne Chain `json:"besu_1,omitempty"`
	BesuTwo Chain `json:"besu_2,omitempty"`
}

type Chain struct {
	URLs []string `json:"URLs,omitempty"`
}
