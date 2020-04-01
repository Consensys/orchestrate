package testutils

import (
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/abi"
	"k8s.io/apimachinery/pkg/util/rand"
)

const contractABI = `[
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "by",
          "type": "uint256"
        }
      ],
      "name": "Incremented",
      "type": "event"
    },
    {
      "constant": false,
      "inputs": [
        {
          "internalType": "uint256",
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "increment",
      "outputs": [],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ]`
const bytecode = "0x6080604052348015600f57600080fd5b5061010a8061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c80637cf5dab014602d575b600080fd5b605660048036036020811015604157600080fd5b81019080803590602001909291905050506058565b005b8060008082825401925050819055507f38ac789ed44572701765277c4d0970f2db1c1a571ed39e84358095ae4eaa54203382604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a15056fea265627a7a7231582047f0937eed22d0cd2ce2eb7c89d1ef3323fa25ba9839d9c7e55fdc20c906396764736f6c63430005100032"
const deployedBytecdode = "0x6080604052348015600f57600080fd5b506004361060285760003560e01c80637cf5dab014602d575b600080fd5b605660048036036020811015604157600080fd5b81019080803590602001909291905050506058565b005b8060008082825401925050819055507f38ac789ed44572701765277c4d0970f2db1c1a571ed39e84358095ae4eaa54203382604051808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018281526020019250505060405180910390a15056fea265627a7a7231582047f0937eed22d0cd2ce2eb7c89d1ef3323fa25ba9839d9c7e55fdc20c906396764736f6c63430005100032"
const eventABI = `{
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "by",
          "type": "uint256"
        }
      ],
      "name": "Incremented",
      "type": "event"
    }
`
const eventSignature = "Incremented(address,uint256)"

// FakeContract returns a new fake contract
func FakeContract() *abi.Contract {
	return &abi.Contract{
		Id: &abi.ContractId{
			Name: rand.String(5),
			Tag:  "v1.0.0",
		},
		Abi:              contractABI,
		Bytecode:         bytecode,
		DeployedBytecode: deployedBytecdode,
	}
}

func FakeEvent() *abi.Event {
	return &abi.Event{
		Signature: eventSignature,
		Abi:       eventABI,
	}
}
