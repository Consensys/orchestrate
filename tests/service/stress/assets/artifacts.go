package assets

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"encoding/json"

	"github.com/consensys/orchestrate/pkg/sdk/client"
	"github.com/consensys/orchestrate/pkg/toolkit/app/log"
	"github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/entities"
	"github.com/consensys/quorum/common/hexutil"
)

var artifactsCtxKey ctxKey = "artifacts"

type Artifact struct {
	Abi              json.RawMessage
	Bytecode         string
	DeployedBytecode string
}

func RegisterNewContract(ctx context.Context, cClient client.ContractClient, artifactPath, name string) (context.Context, error) {
	logger := log.FromContext(ctx).WithField("name", name)
	logger.Debug("registering new contract")
	contract, err := readContract(ctx, artifactPath, fmt.Sprintf("%s.json", name))
	if err != nil {
		return nil, err
	}

	var abi interface{}
	err = json.Unmarshal([]byte(contract.ABI), &abi)
	if err != nil {
		errMsg := "failed to decode contract ABI"
		logger.WithError(err).Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	resp, err := cClient.RegisterContract(ctx, &api.RegisterContractRequest{
		Name:             name,
		Tag:              contract.Tag,
		ABI:              abi,
		Bytecode:         contract.Bytecode,
		DeployedBytecode: contract.DeployedBytecode,
	})

	if err != nil {
		errMsg := "failed to register contract"
		logger.WithError(err).Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	logger.Info("new contract has been registered")
	return contextWithArtifacts(ctx, append(ContextArtifacts(ctx), resp.Name)), nil
}

func readContract(ctx context.Context, artifactsPath, fileName string) (*entities.Contract, error) {
	f, err := os.Open(path.Join(artifactsPath, fileName))
	if err != nil {
		return nil, err
	}

	defer func() {
		err = f.Close()
		if err != nil {
			log.FromContext(ctx).WithError(err).Error("cannot close artifact file")
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
	contract.Bytecode = hexutil.MustDecode(art.Bytecode)
	// Bytecode is an hexstring encoded []byte
	contract.DeployedBytecode = hexutil.MustDecode(art.DeployedBytecode)

	return contract, nil
}

func contextWithArtifacts(ctx context.Context, artifacts []string) context.Context {
	return context.WithValue(ctx, artifactsCtxKey, artifacts)
}

func ContextArtifacts(ctx context.Context) []string {
	v, ok := ctx.Value(artifactsCtxKey).([]string)
	if !ok {
		return []string{}
	}
	return v
}
