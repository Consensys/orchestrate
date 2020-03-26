package utils

type Artifact struct {
	Abi              string `protobuf:"string,2,opt,name=abi,proto3" json:"abi,omitempty"`
	Bytecode         string `protobuf:"string,3,opt,name=bytecode,proto3" json:"bytecode,omitempty"`
	DeployedBytecode string `protobuf:"string,6,opt,name=deployedBytecode,proto3" json:"deployedBytecode,omitempty"`
}
