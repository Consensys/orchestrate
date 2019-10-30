package common

type Artifact struct {
	Abi              []byte `protobuf:"bytes,2,opt,name=abi,proto3" json:"abi,omitempty"`
	Bytecode         []byte `protobuf:"bytes,3,opt,name=bytecode,proto3" json:"bytecode,omitempty"`
	DeployedBytecode []byte `protobuf:"bytes,6,opt,name=deployedBytecode,proto3" json:"deployedBytecode,omitempty"`
}
