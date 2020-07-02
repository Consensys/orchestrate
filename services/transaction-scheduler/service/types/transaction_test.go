package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/utils"
)

func TestTransactionParams_BasicSuccessful(t *testing.T) {
	params := TransactionParams{
		From:            "0x88a5C2d9919e46F883EB62F7b8Dd9d0CC45bc290",
		To:              "0x88a5C2d9919e46F883EB62F7b8Dd9d0CC45bc291",
		MethodSignature: "Constructor()",
	}

	err := utils.GetValidator().Struct(params)
	assert.NoError(t, err)
}

func TestTransactionParams_SuccessfulOneTimeKeyWithoutFrom(t *testing.T) {
	params := TransactionParams{
		To:              "0x88a5C2d9919e46F883EB62F7b8Dd9d0CC45bc291",
		MethodSignature: "Constructor()",
		OneTimeKey:      true,
	}

	err := utils.GetValidator().Struct(params)
	assert.NoError(t, err)
}

func TestTransactionParams_FailWithoutFrom(t *testing.T) {
	params := TransactionParams{
		To:              "0x88a5C2d9919e46F883EB62F7b8Dd9d0CC45bc291",
		MethodSignature: "Constructor()",
	}

	err := utils.GetValidator().Struct(params)
	assert.Error(t, err)
}

func TestDeployContractParams_BasicSuccessful(t *testing.T) {
	params := DeployContractParams{
		From:         "0x88a5C2d9919e46F883EB62F7b8Dd9d0CC45bc290",
		ContractName: "SimpleContract",
	}

	err := utils.GetValidator().Struct(params)
	assert.NoError(t, err)
}

func TestDeployContractParams_SuccessfulOneTimeKeyWithoutFrom(t *testing.T) {
	params := DeployContractParams{
		ContractName: "SimpleContract",
		OneTimeKey:   true,
	}

	err := utils.GetValidator().Struct(params)
	assert.NoError(t, err)
}

func TestDeployContractParams_FailWithoutFrom(t *testing.T) {
	params := DeployContractParams{
		ContractName: "SimpleContract",
	}

	err := utils.GetValidator().Struct(params)
	assert.Error(t, err)
}
