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
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/abi"
)

type artifact struct {
	Abi              json.RawMessage
	Bytecode         string
	DeployedBytecode string
}

type ContractSpec struct {
	Contract *abi.Contract
	JWTToken string
}

func ParseContracts(table *gherkin.PickleStepArgument_PickleTable) ([]*ContractSpec, error) {
	var contractSpecs []*ContractSpec
	headers := table.Rows[0]
	for _, row := range table.Rows[1:] {
		contractSpec := &ContractSpec{Contract: &abi.Contract{}}
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
		raw, err := openArtifact(cell)
		if err != nil {
			return err
		}

		var a artifact
		err = json.Unmarshal(raw, &a)
		if err != nil {
			return err
		}

		// Abi is a UTF-8 encoded string. Therefore, we can make the straightforward transition
		contractSpec.Contract.Abi = string(a.Abi)
		// Bytecode is an hexstring encoded []byte
		contractSpec.Contract.Bytecode = a.Bytecode
		// Bytecode is an hexstring encoded []byte
		contractSpec.Contract.DeployedBytecode = a.DeployedBytecode
	case "name":
		if contractSpec.Contract.Id == nil {
			contractSpec.Contract.Id = &abi.ContractId{}
		}
		contractSpec.Contract.Id.Name = cell
	case "tag":
		if contractSpec.Contract.Id == nil {
			contractSpec.Contract.Id = &abi.ContractId{}
		}
		contractSpec.Contract.Id.Tag = cell
	case "Headers.Authorization":
		contractSpec.JWTToken = cell
	default:
		return fmt.Errorf("unknown field %q", header)
	}
	return nil
}

func openArtifact(fileName string) ([]byte, error) {
	// Loop over all cucumber folders to possibly find file
	// <cucumber_folder>/artifacts/<fileName>
	for _, v := range viper.GetStringSlice("cucumber.paths") {
		f, err := os.Open(path.Join(v, "artifacts", fileName))
		if err != nil {
			continue
		}

		bytes, readErr := ioutil.ReadAll(f)

		err = f.Close()
		if err != nil {
			log.Error(err)
		}

		return bytes, readErr
	}
	return nil, os.ErrNotExist
}
