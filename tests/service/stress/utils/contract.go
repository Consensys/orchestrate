package utils

import (
	"context"
	json2 "encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/containous/traefik/v2/pkg/log"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/abi"
	registry "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/contract-registry/proto"
)

type Artifact struct {
	Abi              json2.RawMessage
	Bytecode         string
	DeployedBytecode string
}

func RegisterNewContract(ctx context.Context, client registry.ContractRegistryClient, artifactPath, name string) error {
	log.FromContext(ctx).Debugf("Registering new contract %s...", name)
	contract, err := readContract(artifactPath, fmt.Sprintf("%s.json", name))
	if err != nil {
		return err
	}

	contract.Id = &abi.ContractId{
		Name: name,
	}
	_, err = client.RegisterContract(ctx, &registry.RegisterContractRequest{
		Contract: contract,
	})

	if err != nil {
		return err
	}

	log.FromContext(ctx).Infof("New contract registered: %s", name)
	return nil
}

func readContract(artifactsPath, fileName string) (*abi.Contract, error) {
	f, err := os.Open(path.Join(artifactsPath, fileName))
	if err != nil {
		return nil, err
	}

	defer func() {
		err = f.Close()
		if err != nil {
			log.WithoutContext().WithError(err).Error("cannot close artifact file")
		}
	}()
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var art Artifact
	err = json.Unmarshal(bytes, &art)
	if err != nil {
		return nil, err
	}

	var contract = &abi.Contract{}
	contract.Abi = string(art.Abi)
	// Bytecode is an hexstring encoded []byte
	contract.Bytecode = art.Bytecode
	// Bytecode is an hexstring encoded []byte
	contract.DeployedBytecode = art.DeployedBytecode

	return contract, nil
}
