package main

import (
	"context"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/ethclient/rpc"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/ethereum/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/ethereum"
)

// This script has been implemented as part of ABC2D hackathon
// in order to e2e test integration of Quorum privateFor transaction in Orchestrate

// We Have been using a simple Meter contract for test :

// pragma solidity ^0.5.2;

// contract Meter  {
//     mapping (address => uint256) public meters;

//     address[] public participants;

//     event Incremented(address indexed participant, uint256 oldV, uint256 newV);

//     function increment(uint256 value) public returns (bool) {
//         require(value > 0);
//         if (meters[msg.sender] == 0) {
//             participants.push(msg.sender);
//         }
//         meters[msg.sender] += value;

//         emit Incremented(msg.sender, meters[msg.sender] - value, meters[msg.sender]);

//         return true;
//     }
// }

// Deployment bytecode: "0x608060405234801561001057600080fd5b506103d8806100206000396000f3fe608060405234801561001057600080fd5b506004361061005e576000357c01000000000000000000000000000000000000000000000000000000009004806335c1d349146100635780637cf5dab0146100d1578063a55d0aaa14610117575b600080fd5b61008f6004803603602081101561007957600080fd5b810190808035906020019092919050505061016f565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6100fd600480360360208110156100e757600080fd5b81019080803590602001909291905050506101ad565b604051808215151515815260200191505060405180910390f35b6101596004803603602081101561012d57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610394565b6040518082815260200191505060405180910390f35b60018181548110151561017e57fe5b906000526020600020016000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600080821115156101bd57600080fd5b60008060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054141561026b5760013390806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505b816000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825401925050819055503373ffffffffffffffffffffffffffffffffffffffff167fcd5ad702c30bb253c9e421ea7f3e00faee62ce859708bfdaf949788e5ba0fdb5836000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054036000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054604051808381526020018281526020019250505060405180910390a260019050919050565b6000602052806000526040600020600091509050548156fea165627a7a72305820f3b866ec54b0b7ae06c5f182fccfd6261283fe48d8714ea60d0a517a8a0f1e120029"
// Increment payload: "0x7cf5dab00000000000000000000000000000000000000000000000000000000000000005"
func main() {
	// Configure logger
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)

	urls := map[string]string{
		"quorum": "http://localhost:22000",
	}

	endpoint := urls["quorum"]

	ethclient.Init(context.Background())

	chainID, _ := ethclient.GlobalClient().Network(context.Background(), endpoint)
	log.Infof("Connected to chain: %v", chainID.Text(16))

	// Create a Envelope for PrivateFor including Cargill & Admin nodes
	e := &envelope.Envelope{
		From: "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4",
		// Call: &common.Call{
		// 	Quorum: &quorum.Quorum{
		// 		Version:     "1.10.1",
		// 		PrivateFrom: "+PXSmQe5XHivlJGJTx7zYp+pv1Om+7M45AE7uzkpdwA=",
		// 		PrivateFor: []string{
		// 			"8uutlYmU5KujtENydTd6VSrNx2WnjH7BGbbKog3DuDo=",
		// 		},
		// 	},
		// },
		Tx: &ethereum.Transaction{
			TxData: &ethereum.TxData{
				// To: "0xddc9ac261ef978e8441de52b0ac11ee7ccaf0aa2",
				Data: "0x608060405234801561001057600080fd5b506103d8806100206000396000f3fe608060405234801561001057600080fd5b506004361061005e576000357c01000000000000000000000000000000000000000000000000000000009004806335c1d349146100635780637cf5dab0146100d1578063a55d0aaa14610117575b600080fd5b61008f6004803603602081101561007957600080fd5b810190808035906020019092919050505061016f565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6100fd600480360360208110156100e757600080fd5b81019080803590602001909291905050506101ad565b604051808215151515815260200191505060405180910390f35b6101596004803603602081101561012d57600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610394565b6040518082815260200191505060405180910390f35b60018181548110151561017e57fe5b906000526020600020016000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600080821115156101bd57600080fd5b60008060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054141561026b5760013390806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505b816000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825401925050819055503373ffffffffffffffffffffffffffffffffffffffff167fcd5ad702c30bb253c9e421ea7f3e00faee62ce859708bfdaf949788e5ba0fdb5836000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054036000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054604051808381526020018281526020019250505060405180910390a260019050919050565b6000602052806000526040600020600091509050548156fea165627a7a72305820f3b866ec54b0b7ae06c5f182fccfd6261283fe48d8714ea60d0a517a8a0f1e120029",
				// Data: "0x7cf5dab00000000000000000000000000000000000000000000000000000000000000005",
				// GasPrice: "0x0",
				Gas:   314454,
				Nonce: 311,
			},
		},
	}

	// Prepare arguments and send transaction
	args := types.Envelope2SendTxArgs(e)
	txHash, err := ethclient.GlobalClient().SendTransaction(context.Background(), endpoint, args)
	if err != nil {
		log.WithError(err).Errorf("Could not send Quorum private for transaction")
		return
	}
	log.Infof("TxHash: %v", txHash.Hex())

	// Wait for receipt
	for {
		receipt, err := ethclient.GlobalClient().TransactionReceipt(context.Background(), endpoint, txHash)
		if receipt == nil {
			log.WithError(err).Warnf("Waiting for receipt")
			time.Sleep(time.Second)
			continue
		}
		r, err := json.Marshal(receipt)
		if err != nil {
			log.WithError(err).Errorf("Could not parse receipt")
		}
		log.Infof("Transaction mined receipt: %q", string(r))
		break
	}
}
