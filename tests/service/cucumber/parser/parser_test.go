package parser

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"

	"github.com/cucumber/godog/gherkin"
)

func TestParser_ParseTxRequest(t *testing.T) {

	testSet := []struct {
		name              string
		rows              [][]*gherkin.TableCell
		expectedEnvelopes []*tx.Envelope
		expectedError     error
	}{
		{
			"test",
			[][]*gherkin.TableCell{
				{
					{Value: "method"},
					{Value: "from"},
					{Value: "to"},
					{Value: "gas"},
					{Value: "gasPrice"},
					{Value: "nonce"},
					{Value: "id"},
					{Value: "data"},
					{Value: "raw"},
					{Value: "txHash"},
					{Value: "chainID"},
					{Value: "chainName"},
					{Value: "chainUUID"},
					{Value: "contractName"},
					{Value: "contractTag"},
					{Value: "methodSignature"},
					{Value: "args"},
					{Value: "privateFor"},
					{Value: "privateFrom"},
					{Value: "privateTxType"},
					{Value: "privacyGroupID"},
					{Value: "headers.foo"},
					{Value: "headers.bar"},
					{Value: "contextLabels.foo"},
					{Value: "contextLabels.bar"},
					{Value: "internalLabels.foo"},
					{Value: "internalLabels.bar"},
				},
				{
					{Value: "ETH_SENDRAWTRANSACTION"},
					{Value: "0x7e654d251da770a068413677967f6d3ea2fea9e4"},
					{Value: "0xe3F5351F8da45aE9150441E3Af21906CCe4cBbc0"},
					{Value: "1"},
					{Value: "2"},
					{Value: "3"},
					{Value: "85b286ac-1353-43c9-b239-3d55d432ab02"},
					{Value: "0xab"},
					{Value: "0xac"},
					{Value: "0xb57324b29aad016267f3a7f2c3f1b59f8bc264e247c41939fa53c55335b67855"},
					{Value: "9"},
					{Value: "testChainName"},
					{Value: "21dafbe2-9dec-49cd-96f0-296c8adc0f41"},
					{Value: "testContractName"},
					{Value: "testContractTag"},
					{Value: "constructor(int64)"},
					{Value: "10"},
					{Value: "dGVzdA=="},
					{Value: "UHJpdmF0ZUZyb20="},
					{Value: "testPrivateTxType"},
					{Value: "testPrivacyGroupID"},
					{Value: "fooValueHeader"},
					{Value: "barValueHeader"},
					{Value: "fooValueContextLabels"},
					{Value: "barValueContextLabels"},
					{Value: "fooValueInternalLabels"},
					{Value: "barValueInternalLabels"},
				},
			},
			[]*tx.Envelope{
				tx.NewEnvelope().
					SetMethod(tx.Method_ETH_SENDRAWTRANSACTION).
					MustSetFromString("0x7e654d251da770a068413677967f6d3ea2fea9e4").
					MustSetToString("0xe3F5351F8da45aE9150441E3Af21906CCe4cBbc0").
					SetGas(1).
					SetGasPrice(big.NewInt(2)).
					SetNonce(3).
					SetID("85b286ac-1353-43c9-b239-3d55d432ab02").
					MustSetDataString("0xab").
					MustSetRawString("0xac").
					MustSetTxHashString("0xb57324b29aad016267f3a7f2c3f1b59f8bc264e247c41939fa53c55335b67855").
					SetChainID(big.NewInt(9)).
					SetChainName("testChainName").
					SetChainUUID("21dafbe2-9dec-49cd-96f0-296c8adc0f41").
					SetContractName("testContractName").
					SetContractTag("testContractTag").
					SetMethodSignature("constructor(int64)").
					SetArgs([]string{"10"}).
					SetPrivateFor([]string{"dGVzdA=="}).
					SetPrivateFrom("UHJpdmF0ZUZyb20=").
					SetPrivateTxType("testPrivateTxType").
					SetPrivacyGroupID("testPrivacyGroupID").
					SetHeadersValue("foo", "fooValueHeader").
					SetHeadersValue("bar", "barValueHeader").
					SetContextLabelsValue("foo", "fooValueContextLabels").
					SetContextLabelsValue("bar", "barValueContextLabels").
					SetInternalLabelsValue("foo", "fooValueInternalLabels").
					SetInternalLabelsValue("bar", "barValueInternalLabels"),
			},
			nil,
		},
	}

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			parser := New()

			rows := make([]*gherkin.TableRow, 0)
			for _, row := range test.rows {
				rows = append(rows, &gherkin.TableRow{Cells: row})
			}

			evlps, err := parser.ParseEnvelopes("test-1", &gherkin.DataTable{Rows: rows})
			assert.Equal(t, test.expectedError, err)

			for i := range evlps {
				assert.True(t, reflect.DeepEqual(test.expectedEnvelopes[i], evlps[i]))
			}
		})
	}
}
