package txschedulertypes

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
		Annotations:     Annotations{OneTimeKey: true},
	}

	err := utils.GetValidator().Struct(params)
	assert.NoError(t, err)
}

func TestTransactionParams_FailWithoutFrom(t *testing.T) {
	params := TransactionParams{
		To:              "0x88a5C2d9919e46F883EB62F7b8Dd9d0CC45bc291",
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
				Annotations: Annotations{
					RetryPolicy: GasPriceRetryParams{
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
				Annotations: Annotations{
					RetryPolicy: GasPriceRetryParams{
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
				Value: "error",
			},
			true,
		},

		{
			"Retry params retry",
			&TransferParams{
				RetryPolicy: GasPriceRetryParams{
					Limit: 0,
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
		Annotations:  Annotations{OneTimeKey: true},
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
		ContractName: "SimpleContract",
		Annotations:  Annotations{Priority: "invalidPriority"},
	}

	err := params.Validate()
	assert.Error(t, err)
}

func TestRetryParams_Validation(t *testing.T) {
	testSet := []struct {
		name          string
		params        GasPriceRetryParams
		expectedError bool
	}{
		{
			"Limit not filled if Increment is filled",
			GasPriceRetryParams{
				Increment: 1.1,
			},
			true,
		},
		{
			"Limit not filled if IncrementLevel is filled",
			GasPriceRetryParams{
				IncrementLevel: "low",
			},
			true,
		},
		{
			"Increment or IncrementLevel not filled if Limit is filled",
			GasPriceRetryParams{
				Limit: 1.1,
			},
			true,
		},
		{
			"No error all fields are filled with Increment",
			GasPriceRetryParams{
				Interval:  "1m",
				Increment: 1.1,
				Limit:     1.2,
			},
			false,
		},
		{
			"No error all fields are filled with IncrementLevel",
			GasPriceRetryParams{
				Interval:       "1m",
				IncrementLevel: "medium",
				Limit:          1.2,
			},
			false,
		},
		{
			"Interval is not a duration",
			GasPriceRetryParams{
				Interval:  "1_m",
				Increment: 1.1,
				Limit:     1.2,
			},
			true,
		},
		{
			"Increment > Limit",
			GasPriceRetryParams{
				Interval:  "1m",
				Increment: 1.3,
				Limit:     1.2,
			},
			true,
		},
		{
			"invalid IncrementLevel",
			GasPriceRetryParams{
				Interval:       "1m",
				IncrementLevel: "l0w",
				Limit:          1.2,
			},
			true,
		},
		{
			"mutual exclusion between Increment and IncrementLevel",
			GasPriceRetryParams{
				Interval:       "1m",
				IncrementLevel: utils.GasIncrementMedium,
				Increment:      1.3,
				Limit:          1.2,
			},
			true,
		},
		{
			"No error when Increment = Limit",
			GasPriceRetryParams{
				Interval:  "1m",
				Increment: 1.1,
				Limit:     1.1,
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
