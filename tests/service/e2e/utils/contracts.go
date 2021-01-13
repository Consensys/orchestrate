package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	gherkin "github.com/cucumber/messages-go/v10"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/entities"
)

type Artifact struct {
	Abi              json.RawMessage
	Bytecode         string
	DeployedBytecode string
}

type ContractSpec struct {
	Contract *entities.Contract
	JWTToken string
}

func ParseContracts(table *gherkin.PickleStepArgument_PickleTable) ([]*ContractSpec, error) {
	var contractSpecs []*ContractSpec
	headers := table.Rows[0]
	for _, row := range table.Rows[1:] {
		contractSpec := &ContractSpec{Contract: &entities.Contract{}}
		err := ParseContract(headers, row, contractSpec)
		if err != nil {
			return nil, err
		}
		contractSpecs = append(contractSpecs, contractSpec)
	}
	return contractSpecs, nil
}

func ParseContract(headers, row *gherkin.PickleStepArgument_PickleTable_PickleTableRow, contractSpec *ContractSpec) error {
	for i, cell := range row.Cells {
		err := ParseContractCell(headers.Cells[i].Value, cell.Value, contractSpec)
		if err != nil {
			return err
		}
	}
	return nil
}

func ParseContractCell(header, cell string, contractSpec *ContractSpec) error {
	switch header {
	case "artifacts":
		raw, err := OpenArtifact(cell, viper.GetString("artifacts.path"))
		if err != nil {
			return err
		}

		var a Artifact
		err = json.Unmarshal(raw, &a)
		if err != nil {
			return err
		}

		// Abi is a UTF-8 encoded string. Therefore, we can make the straightforward transition
		contractSpec.Contract.ABI = string(a.Abi)
		// Bytecode is an hexstring encoded []byte
		contractSpec.Contract.Bytecode = a.Bytecode
		// Bytecode is an hexstring encoded []byte
		contractSpec.Contract.DeployedBytecode = a.DeployedBytecode
	case "name":
		contractSpec.Contract.Name = cell
	case "tag":
		contractSpec.Contract.Tag = cell
	case "Headers.Authorization":
		contractSpec.JWTToken = cell
	default:
		return fmt.Errorf("unknown field %q", header)
	}
	return nil
}

func OpenArtifact(fileName, artifactPath string) ([]byte, error) {
	f, err := os.Open(path.Join(artifactPath, fileName))
	if err != nil {
		return nil, err
	}

	bytes, readErr := ioutil.ReadAll(f)

	err = f.Close()
	if err != nil {
		log.Error(err)
	}

	return bytes, readErr
}
