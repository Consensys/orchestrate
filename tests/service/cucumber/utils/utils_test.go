// +build unit

package utils

import (
	"math/big"
	"reflect"
	"testing"

	gherkin "github.com/cucumber/messages-go/v10"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/tx"
)

type TestStruct struct {
	TestString      string
	TestInt8        int8
	TestUint64      uint64
	TestBool        bool
	TestSliceString []string
	TestSliceInt    []int
	TestSliceUint   []uint
	TestSliceBool   []bool
	NestedStruct    *TestStruct
}

func TestParseTable(t *testing.T) {
	testSet := []struct {
		name                    string
		inputInterface          interface{}
		inputTable              *gherkin.PickleStepArgument_PickleTable
		expectedInterfaceSlices []interface{}
		expectedError           error
	}{
		{
			"single slice parse",
			TestStruct{},
			&gherkin.PickleStepArgument_PickleTable{
				Rows: []*gherkin.PickleStepArgument_PickleTable_PickleTableRow{
					{
						Cells: []*gherkin.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
							{Value: "testString"},
							{Value: "testInt8"},
							{Value: "testUint64"},
							{Value: "testBool"},
							{Value: "testSliceString"},
							{Value: "nestedStruct.testBool"},
						},
					},
					{
						Cells: []*gherkin.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
							{Value: "stepTable"},
							{Value: "-2"},
							{Value: "3"},
							{Value: "true"},
							{Value: "[\"str1\",\"str2\",\"str3\"]"},
							{Value: "true"},
						},
					},
				},
			},
			[]interface{}{
				&TestStruct{
					TestString:      "stepTable",
					TestInt8:        -2,
					TestUint64:      3,
					TestBool:        true,
					TestSliceString: []string{"str1", "str2", "str3"},
					NestedStruct:    &TestStruct{TestBool: true},
				},
			},
			nil,
		},
		{
			"multiple slice parse",
			TestStruct{},
			&gherkin.PickleStepArgument_PickleTable{
				Rows: []*gherkin.PickleStepArgument_PickleTable_PickleTableRow{
					{
						Cells: []*gherkin.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
							{Value: "testString"},
							{Value: "testInt8"},
							{Value: "testUint64"},
							{Value: "testBool"},
							{Value: "nestedStruct.testBool"},
							{Value: "nestedStruct.testSliceUint"},
						},
					},
					{
						Cells: []*gherkin.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
							{Value: "stepTable"},
							{Value: "-2"},
							{Value: "3"},
							{Value: "true"},
							{Value: "true"},
							{Value: "[5,9,1]"},
						},
					},
					{
						Cells: []*gherkin.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
							{Value: "stepTable"},
							{Value: "-2"},
							{Value: "3"},
							{Value: "true"},
							{Value: "true"},
							{Value: "[4,9,1]"},
						},
					},
				},
			},
			[]interface{}{
				&TestStruct{
					TestString:   "stepTable",
					TestInt8:     -2,
					TestUint64:   3,
					TestBool:     true,
					NestedStruct: &TestStruct{TestBool: true, TestSliceUint: []uint{5, 9, 1}},
				},
				&TestStruct{
					TestString:   "stepTable",
					TestInt8:     -2,
					TestUint64:   3,
					TestBool:     true,
					NestedStruct: &TestStruct{TestBool: true, TestSliceUint: []uint{4, 9, 1}},
				},
			},
			nil,
		},
		{
			"nested struct",
			TestStruct{},
			&gherkin.PickleStepArgument_PickleTable{
				Rows: []*gherkin.PickleStepArgument_PickleTable_PickleTableRow{
					{
						Cells: []*gherkin.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
							{Value: "nestedStruct.testBool"},
							{Value: "nestedStruct.nestedStruct.nestedStruct.testSliceBool"},
						},
					},
					{
						Cells: []*gherkin.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
							{Value: "true"},
							{Value: "[false,false,true]"},
						},
					},
				},
			},
			[]interface{}{
				&TestStruct{
					NestedStruct: &TestStruct{
						TestBool: true,
						NestedStruct: &TestStruct{
							NestedStruct: &TestStruct{
								TestSliceBool: []bool{false, false, true},
							},
						},
					},
				},
			},
			nil,
		},
	}

	for _, test := range testSet {
		test := test
		t.Run(test.name, func(t *testing.T) {
			interfaceSlices, err := ParseTable(test.inputInterface, test.inputTable)
			assert.Equal(t, test.expectedError, err)
			assert.True(t, reflect.DeepEqual(interfaceSlices, test.expectedInterfaceSlices))
		})
	}
}

func TestParseEnvelope(t *testing.T) {

	testSet := []struct {
		name              string
		rows              [][]*gherkin.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell
		expectedEnvelopes []*tx.Envelope
		expectedError     error
	}{
		{
			"test",
			[][]*gherkin.PickleStepArgument_PickleTable_PickleTableRow_PickleTableCell{
				{
					{Value: "Method"},
					{Value: "From"},
					{Value: "To"},
					{Value: "Gas"},
					{Value: "GasPrice"},
					{Value: "Nonce"},
					{Value: "ID"},
					{Value: "Data"},
					{Value: "Raw"},
					{Value: "TxHash"},
					{Value: "ChainID"},
					{Value: "ChainName"},
					{Value: "ChainUUID"},
					{Value: "ContractName"},
					{Value: "ContractTag"},
					{Value: "MethodSignature"},
					{Value: "Args"},
					{Value: "PrivateFor"},
					{Value: "PrivateFrom"},
					{Value: "PrivateTxType"},
					{Value: "PrivacyGroupID"},
					{Value: "Headers.foo"},
					{Value: "Headers.bar"},
					{Value: "ContextLabels.foo"},
					{Value: "ContextLabels.bar"},
					{Value: "InternalLabels.foo"},
					{Value: "InternalLabels.bar"},
				},
				{
					{Value: "0"},
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
					{Value: "[\"10\"]"},
					{Value: "[\"dGVzdA==\"]"},
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
					SetJobType(tx.JobType_ETH_TX).
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

			rows := make([]*gherkin.PickleStepArgument_PickleTable_PickleTableRow, 0)
			for _, row := range test.rows {
				rows = append(rows, &gherkin.PickleStepArgument_PickleTable_PickleTableRow{Cells: row})
			}

			evlps, err := ParseEnvelope(&gherkin.PickleStepArgument_PickleTable{Rows: rows})
			assert.Equal(t, test.expectedError, err)

			for i := range evlps {
				assert.True(t, reflect.DeepEqual(test.expectedEnvelopes[i], evlps[i]))
			}
		})
	}
}
