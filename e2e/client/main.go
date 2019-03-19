package main

import (
	"context"
	"math/big"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/golang/protobuf/ptypes"
	log "github.com/sirupsen/logrus"
	common "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
	store "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/context-store"
	ethereum "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/ethereum"
	trace "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.DialContext(
		context.Background(),
		"localhost:8080",
		grpc.WithInsecure(),
	)
	if err != nil {
		panic(err)
	}
	client := store.NewStoreClient(conn)

	txData := (&ethereum.TxData{}).
		SetNonce(10).
		SetTo(ethcommon.HexToAddress("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C")).
		SetValue(big.NewInt(100000)).
		SetGas(2000).
		SetGasPrice(big.NewInt(200000)).
		SetData(hexutil.MustDecode("0xabcd"))

	tr := &trace.Trace{
		Chain:    &common.Chain{Id: "0x3"},
		Metadata: &trace.Metadata{Id: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11"},
		Tx: &ethereum.Transaction{
			TxData: txData,
			Raw:    "0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80",
			Hash:   "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210",
		},
	}

	resp, err := client.Store(context.Background(), &store.StoreRequest{
		Trace: tr,
	})
	if err != nil {
		log.WithError(err).Errorf("Could not store")
	}

	timestamp, err := ptypes.Timestamp(resp.LastUpdated)
	if err != nil {
		log.WithError(err).Errorf("Could not store")
	}

	log.WithFields(log.Fields{
		"status": resp.Status,
		"at":     timestamp,
	}).Infof("Trace stored")
	conn.Close()
}
