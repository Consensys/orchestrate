package main

import (
	"context"
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	ethclient "gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/ethclient"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/ethereum.git/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
)

// This script has been implemented as part of ABC2D hackathon
// in order to e2e test integration of Quorum privateFor transaction in Core-Stack

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

// Staging environment set-up:

// Node 1: Admin-Staging
// https://e0bowvhhuf:UVmqpGBOOkQbKDacmLQsaL0QfD0-EuLfpNzP4or9m8k@e0w1quks6d-e0lospwmr9-rpc.eu-central-1.kaleido.io
// 8uutlYmU5KujtENydTd6VSrNx2WnjH7BGbbKog3DuDo=
// 0x5dD217f9E9009871Dc41695fBd942076EFe15f3C

// Node 2: Bunge-Staging
// https://e0k0h01n0w:Cf9XOEgS_zBnKmKoTw1NYDccHNiDLeSCpEFmFcb36hU@e0w1quks6d-e0pm3msflc-rpc.eu-central-1.kaleido.io
// FyJZPKRSd8Bsqu0V+qXpx/Zm986PcldsTXMk8QxHCWQ=
// 0x566D44FB1442d206439F23E09885Ee365Bb1E43f

// Node 3: Cargill-Staging
// https://e0jcrldvk7:eOjDFyulgTHhUm1C5hAjvW9fAPZBM-eciSfA0bzlkyU@e0w1quks6d-e0x1svr9gn-rpc.eu-central-1.kaleido.io
// +PXSmQe5XHivlJGJTx7zYp+pv1Om+7M45AE7uzkpdwA=
// 0x22460fa1b318897934fF1bb3dfeA19Ed9B218dB4
func main() {
	// Configure logger
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)

	// Initialize client
	viper.Set("eth.clients", []string{
		"https://e0jcrldvk7:eOjDFyulgTHhUm1C5hAjvW9fAPZBM-eciSfA0bzlkyU@e0w1quks6d-e0x1svr9gn-rpc.eu-central-1.kaleido.io",
	})

	ethclient.Init(context.Background())

	chain := ethclient.GlobalClient().Networks(context.Background())[0]
	log.Infof("Connected to chain: %v", chain.Text(16))

	// Create a Envelope for PrivateFor including Cargill & Admin nodes
	e := &envelope.Envelope{
		Sender: &common.Account{
			Addr: "0x22460fa1b318897934fF1bb3dfeA19Ed9B218dB4",
		},
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

	// Prepare arguments and send trasaction
	args := types.Envelope2SendTxArgs(e)
	txHash, err := ethclient.GlobalClient().SendTransaction(context.Background(), chain, args)
	if err != nil {
		log.WithError(err).Errorf("Could not send Quorum private for transaction")
		return
	}
	log.Infof("TxHash: %v", txHash.Hex())

	// Wait for receipt
	for {
		receipt, err := ethclient.GlobalClient().TransactionReceipt(context.Background(), chain, txHash)
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
