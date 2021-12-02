package api

import (
	"testing"

	"github.com/consensys/orchestrate/pkg/utils"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestTransactionParams_BasicSuccessful(t *testing.T) {
	from := ethcommon.HexToAddress("0x88a5C2d9919e46F883EB62F7b8Dd9d0CC45bc290")
	params := TransactionParams{
		From:            &from,
		To:              utils.ToPtr(ethcommon.HexToAddress("0x88a5C2d9919e46F883EB62F7b8Dd9d0CC45bc290")).(*ethcommon.Address),
		MethodSignature: "Constructor()",
	}

	err := params.Validate()
	assert.NoError(t, err)
}

func TestTransactionParams_SuccessfulOneTimeKeyWithoutFrom(t *testing.T) {
	params := TransactionParams{
		To:              utils.ToPtr(ethcommon.HexToAddress("0x88a5C2d9919e46F883EB62F7b8Dd9d0CC45bc290")).(*ethcommon.Address),
		MethodSignature: "Constructor()",
		OneTimeKey:      true,
	}

	err := utils.GetValidator().Struct(params)
	assert.NoError(t, err)
}

func TestTransactionParams_FailWithoutFrom(t *testing.T) {
	params := TransactionParams{
		To:              utils.ToPtr(ethcommon.HexToAddress("0x88a5C2d9919e46F883EB62F7b8Dd9d0CC45bc291")).(*ethcommon.Address),
		MethodSignature: "Constructor()",
	}

	err := params.Validate()
	assert.Error(t, err)
}

func TestTransactionParams_Validation(t *testing.T) {
	testSet := []struct {
		name          string
		params        *TransactionParams
		expectedError bool
	}{
		{
			"Validator error",
			&TransactionParams{
				Value: nil,
			},
			true,
		},
		{
			"PrivateParams error",
			&TransactionParams{
				PrivateFor:     []string{"test"},
				PrivacyGroupID: "test",
			},
			true,
		},
		{
			"Retry params retry",
			&TransactionParams{
				GasPricePolicy: GasPriceParams{
					RetryPolicy: RetryParams{
						Limit: 0,
					},
				},
			},
			true,
		},
	}

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.params.Validate()
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeployContractParams_Validation(t *testing.T) {
	testSet := []struct {
		name          string
		params        *DeployContractParams
		expectedError bool
	}{
		{
			"Validator error",
			&DeployContractParams{
				Value: nil,
			},
			true,
		},
		{
			"PrivateParams error",
			&DeployContractParams{
				PrivateFor:     []string{"test"},
				PrivacyGroupID: "test",
			},
			true,
		},
		{
			"Retry params retry",
			&DeployContractParams{
				GasPricePolicy: GasPriceParams{
					RetryPolicy: RetryParams{
						Limit: 0,
					},
				},
			},
			true,
		},
	}

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.params.Validate()
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTransferParams_Validation(t *testing.T) {
	testSet := []struct {
		name          string
		params        *TransferParams
		expectedError bool
	}{
		{
			"Validator error",
			&TransferParams{
				Value: nil,
			},
			true,
		},

		{
			"Retry params retry",
			&TransferParams{
				GasPricePolicy: GasPriceParams{
					RetryPolicy: RetryParams{
						Limit: 0,
					},
				},
			},
			true,
		},
	}

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.params.Validate()
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeployContractParams_BasicSuccessful(t *testing.T) {
	from := ethcommon.HexToAddress("0x88a5C2d9919e46F883EB62F7b8Dd9d0CC45bc290")
	params := DeployContractParams{
		From:         &from,
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

	err := params.Validate()
	assert.Error(t, err)
}

func TestParams_Priority(t *testing.T) {
	params := DeployContractParams{
		ContractName:   "SimpleContract",
		GasPricePolicy: GasPriceParams{Priority: "invalidPriority"},
	}

	err := params.Validate()
	assert.Error(t, err)
}

func TestRetryParams_Validation(t *testing.T) {
	testSet := []struct {
		name          string
		params        RetryParams
		expectedError bool
	}{
		{
			"Limit not filled if Increment is filled",
			RetryParams{
				Increment: 1.1,
			},
			true,
		},
		{
			"Increment not filled if Limit is filled",
			RetryParams{
				Limit: 1.1,
			},
			true,
		},
		{
			"No error all fields are filled with Increment",
			RetryParams{
				Interval:  "1m",
				Increment: 1.1,
				Limit:     1.2,
			},
			false,
		},
		{
			"No error all fields are filled with Increment",
			RetryParams{
				Interval:  "1s",
				Increment: 1.1,
				Limit:     1.2,
			},
			false,
		},
		{
			"Interval is not a duration",
			RetryParams{
				Interval:  "1_m",
				Increment: 1.1,
				Limit:     1.2,
			},
			true,
		},
		{
			"Interval duration too low",
			RetryParams{
				Interval:  "100ms",
				Increment: 1.1,
				Limit:     1.2,
			},
			true,
		},
		{
			"Amount of retries exceeds limit",
			RetryParams{
				Interval:  "1s",
				Increment: 0.1,
				Limit:     1.2,
			},
			true,
		},
	}

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.params.Validate()
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
