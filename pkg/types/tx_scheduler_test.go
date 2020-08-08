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

	err := params.Validate()
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

func TestTransactionParams_Validation(t *testing.T) {
	testSet := []struct {
		name          string
		params        *TransactionParams
		expectedError bool
	}{
		{
			"Validator error",
			&TransactionParams{
				Value: "error",
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
				Retry: &GasPriceRetryParams{
					GasPriceLimit: 0,
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
				Value: "error",
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
				Retry: &GasPriceRetryParams{
					GasPriceLimit: 0,
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
				Value: "error",
			},
			true,
		},

		{
			"Retry params retry",
			&TransferParams{
				Retry: &GasPriceRetryParams{
					GasPriceLimit: 0,
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

func TestParams_Priority(t *testing.T) {
	params := DeployContractParams{
		ContractName: "SimpleContract",
		Priority:     "invalidPriority",
	}

	err := utils.GetValidator().Struct(params)
	assert.Error(t, err)
}

func TestRetryParams_Validation(t *testing.T) {
	testSet := []struct {
		name          string
		params        GasPriceRetryParams
		expectedError bool
	}{
		{
			"Error GasPriceIncrement, GasPriceLimit not filled",
			GasPriceRetryParams{
				Interval: "1m",
			},
			true,
		},
		{
			"Error GasPriceLimit not filled",
			GasPriceRetryParams{
				Interval:          "1m",
				GasPriceIncrement: 1.1,
			},
			true,
		},
		{
			"Error GasPriceIncrement or GasPriceIncrementLevel not filled",
			GasPriceRetryParams{
				Interval:      "1m",
				GasPriceLimit: 1.1,
			},
			true,
		},
		{
			"Error Interval, GasPriceIncrement not filled",
			GasPriceRetryParams{
				GasPriceLimit: 1.2,
			},
			true,
		},
		{
			"Error Interval not filled",
			GasPriceRetryParams{
				GasPriceIncrementLevel: "low",
				GasPriceLimit:          1.2,
			},
			true,
		},
		{
			"Error Interval and GasPriceLimit not filled",
			GasPriceRetryParams{
				GasPriceIncrementLevel: "low",
			},
			true,
		},
		{
			"Error Interval, GasPriceLimit not filled",
			GasPriceRetryParams{
				GasPriceIncrement: 1.2,
			},
			true,
		},
		{
			"No error all fields are filled",
			GasPriceRetryParams{
				Interval:          "1m",
				GasPriceIncrement: 1.1,
				GasPriceLimit:     1.2,
			},
			false,
		},
		{
			"No error all fields are filled",
			GasPriceRetryParams{
				Interval:               "1m",
				GasPriceIncrementLevel: "medium",
				GasPriceLimit:          1.2,
			},
			false,
		},
		{
			"Error Interval is not a duration",
			GasPriceRetryParams{
				Interval:          "1_m",
				GasPriceIncrement: 1.1,
				GasPriceLimit:     1.2,
			},
			true,
		},
		{
			"Error GasPriceIncrement > GasPriceLimit",
			GasPriceRetryParams{
				Interval:          "1m",
				GasPriceIncrement: 1.3,
				GasPriceLimit:     1.2,
			},
			true,
		},
		{
			"Error invalid GasPriceIncrementLevel",
			GasPriceRetryParams{
				Interval:               "1m",
				GasPriceIncrementLevel: "l0w",
				GasPriceLimit:          1.2,
			},
			true,
		},
		{
			"Error mutual exclusion between GasPriceIncrement and GasPriceIncrementLevel",
			GasPriceRetryParams{
				Interval:               "1m",
				GasPriceIncrementLevel: "low",
				GasPriceIncrement:      1.3,
				GasPriceLimit:          1.2,
			},
			true,
		},
		{
			"No error when GasPriceIncrement = GasPriceLimit",
			GasPriceRetryParams{
				Interval:          "1m",
				GasPriceIncrement: 1.1,
				GasPriceLimit:     1.1,
			},
			false,
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
