package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	evlpstore "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/services_tmp/envelope-store"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/envelope"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/envelope-store.git/client"
)

func main() {
	client.Init(context.Background())
	defer client.Close()

	// Store envelope
	evlp := &envelope.Envelope{
		Chain:    chain.FromInt(888),
		Metadata: &envelope.Metadata{Id: "6be0-bc19-900b-1ef8-bb6d-61b9-ad38-ba12"},
		Tx: &ethereum.Transaction{
			Raw:  ethereum.HexToData("0xf86c0184ee6b2800a2529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80"),
			Hash: ethereum.HexToHash("0x6a0caf026cb1f012abe19e9e02c53f23713b0033d7a72e534136104b5447a21a"),
		},
	}
	resp, err := client.GlobalEnvelopeStoreClient().Store(
		context.Background(),
		&evlpstore.StoreRequest{
			Envelope: evlp,
		})
	if err != nil {
		log.WithError(err).Errorf("Could not store")
	}

	log.WithFields(log.Fields{
		"status": resp.GetStatusInfo().GetStatus().String(),
		"at":     resp.GetStatusInfo().StoredAtTime(),
	}).Infof("Envelope stored")

	res, err := client.GlobalEnvelopeStoreClient().LoadByTxHash(
		context.Background(),
		&evlpstore.LoadByTxHashRequest{
			Chain:  evlp.GetChain(),
			TxHash: evlp.GetTx().GetHash(),
		})
	if err != nil {
		log.WithError(err).Errorf("Could not load envelope")
	}
	log.WithFields(log.Fields{
		"status":   res.GetStatusInfo().GetStatus().String(),
		"chain.id": res.GetEnvelope().GetChain().ID().String(),
	}).Infof("Envelope loaded")
}
