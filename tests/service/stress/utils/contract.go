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
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/sdk/client"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/api"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
)

type Artifact struct {
	Abi              json2.RawMessage
	Bytecode         string
	DeployedBytecode string
}

func RegisterNewContract(ctx context.Context, cClient client.ContractClient, artifactPath, name string) error {
	log.FromContext(ctx).Debugf("Registering new contract %s...", name)
	contract, err := readContract(artifactPath, fmt.Sprintf("%s.json", name))
	if err != nil {
		return err
	}

	var abi interface{}
	err = json.Unmarshal([]byte(contract.ABI), &abi)
	if err != nil {
		return err
	}

	_, err = cClient.RegisterContract(ctx, &api.RegisterContractRequest{
		Name:             name,
		Tag:              contract.Tag,
		ABI:              abi,
		Bytecode:         contract.Bytecode,
		DeployedBytecode: contract.DeployedBytecode,
	})

	if err != nil {
		return err
	}

	log.FromContext(ctx).Infof("New contract registered: %s", name)
	return nil
}

func readContract(artifactsPath, fileName string) (*entities.Contract, error) {
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

	var contract = &entities.Contract{}
	contract.ABI = string(art.Abi)
	// Bytecode is an hexstring encoded []byte
	contract.Bytecode = art.Bytecode
	// Bytecode is an hexstring encoded []byte
	contract.DeployedBytecode = art.DeployedBytecode

	return contract, nil
}
