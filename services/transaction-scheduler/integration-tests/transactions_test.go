// +build integration

package integrationtests

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/transaction-scheduler/transaction-scheduler/types"
	"net/http"
	"testing"
)

// TransactionsTestSuite is a test suite for Transaction API jobs controller
type TransactionsTestSuite struct {
	suite.Suite
	baseURL string
	env     *IntegrationEnvironment
}

func (s *TransactionsTestSuite) TestTransactions_Validation() {
	s.T().Run("should fail with 400 if payload is invalid", func(t *testing.T) {
		resp := s.sendTransaction(&types.TransactionRequest{
			BaseTransactionRequest: types.BaseTransactionRequest{
				IdempotencyKey: "myID",
				ChainID:        "chainID",
			},
			Params: types.TransactionParams{
				From:            "0x7E654d251Da770A068413677967F6d3Ea2FeA9E4",
				To:              "0x905B88EFf8Bda1543d4d6f4aA05afef143D27E18",
				MethodSignature: "constructor()",
			},
		})

		assert.Equal(t, 400, resp.StatusCode)
	})
}

func (s *TransactionsTestSuite) sendTransaction(txRequest *types.TransactionRequest) *http.Response {
	request, err := json.Marshal(txRequest)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(s.baseURL+"/send", "application/json", bytes.NewBuffer(request))
	if err != nil {
		panic(err)
	}

	return resp
}
