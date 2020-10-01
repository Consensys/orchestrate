// +build unit

package tx

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/abi"
	ierror "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/error"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/ethereum"
)

func TestNewEnvelope(t *testing.T) {
	b := NewEnvelope()
	assert.NotNil(t, b, "Should not be nil")
	assert.NotNil(t, b.GetHeaders(), "Should not be nil")
	assert.NotNil(t, b.GetContextLabels(), "Should not be nil")
	assert.NotNil(t, b.GetErrors(), "Should not be nil")
	assert.NotNil(t, b.GetInternalLabels(), "Should not be nil")
}

func TestEnvelope_SetID(t *testing.T) {
	id := uuid.Must(uuid.NewV4()).String()
	b := NewEnvelope().SetID(id)
	assert.Equal(t, id, b.GetID(), "Should be equal")
}

func TestEnvelope_SetJobUUID(t *testing.T) {
	id := uuid.Must(uuid.NewV4()).String()
	b := NewEnvelope().SetJobUUID(id)
	assert.Equal(t, id, b.GetJobUUID(), "Should be equal")
}

func TestEnvelope_GetErrors(t *testing.T) {
	testError := errors.FromError(fmt.Errorf("test"))
	b := NewEnvelope().AppendError(testError)
	assert.True(t, reflect.DeepEqual(b.GetErrors(), []*ierror.Error{testError}), "Should be equal")
}

func TestEnvelope_Error(t *testing.T) {
	b := NewEnvelope()
	assert.Empty(t, b.Error(), "Should be equal")

	testError := errors.FromError(fmt.Errorf("test"))
	b = NewEnvelope().AppendError(testError)
	assert.Equal(t, "[\"FF000@: test\"]", b.Error(), "Should be equal")
}

func TestEnvelope_AppendError(t *testing.T) {
	testError := errors.FromError(fmt.Errorf("test"))
	b := NewEnvelope().AppendError(testError)
	assert.True(t, reflect.DeepEqual(b.Errors, []*ierror.Error{testError}), "Should be equal")
}

func TestEnvelope_AppendErrors(t *testing.T) {
	testErrors := []*ierror.Error{errors.FromError(fmt.Errorf("test1")), errors.FromError(fmt.Errorf("test2"))}
	b := NewEnvelope().AppendErrors(testErrors)
	assert.True(t, reflect.DeepEqual(b.Errors, testErrors), "Should be equal")
}

func TestEnvelope_SetReceipt(t *testing.T) {
	receipt := &ethereum.Receipt{TxHash: "test"}
	b := NewEnvelope().SetReceipt(receipt)
	assert.True(t, reflect.DeepEqual(b.GetReceipt(), receipt), "Should be equal")

}

func TestEnvelope_Carrier(t *testing.T) {
	b := NewEnvelope()
	assert.NotNil(t, b.Carrier(), "Should be equal")
}

func TestEnvelope_OnlyWarnings(t *testing.T) {
	testError := errors.FromError(fmt.Errorf("test"))
	b := NewEnvelope().AppendError(testError)
	assert.False(t, b.OnlyWarnings(), "Should be equal")

	testWarning := errors.Warningf("test")
	b = NewEnvelope().AppendError(testWarning)
	assert.True(t, b.OnlyWarnings(), "Should be equal")
}

func TestEnvelope_Headers(t *testing.T) {
	b := NewEnvelope().SetHeadersValue("key", "value")
	assert.Equal(t, "value", b.GetHeadersValue("key"), "Should be equal")
}

func TestEnvelope_ContextLabels(t *testing.T) {
	b := NewEnvelope().SetContextLabelsValue("key", "value")
	assert.Equal(t, "value", b.GetContextLabelsValue("key"), "Should be equal")

	_ = b.SetContextLabels(map[string]string{"newLabelKey": "newLabelValue"})
	assert.Equal(t, "newLabelValue", b.GetContextLabelsValue("newLabelKey"), "Should be equal")
}

func TestEnvelope_Validate(t *testing.T) {
	b := NewEnvelope().
		SetID("jobUUID").
		SetHeadersValue("testHeaderKey", "testHeaderValue").
		SetChainName("chainName").
		MustSetFromString("0x1").
		SetJobType(JobType_ETH_TX).
		MustSetToString("0x2").
		SetGas(11).
		SetGasPrice(big.NewInt(12)).
		SetValue(big.NewInt(13)).
		SetNonce(14).
		SetData([]byte{1}).
		SetContractName("testContractName").
		SetContractTag("testContractTag").
		SetMethodSignature("testMethodSignature").
		SetArgs([]string{"testArg1", "testArg2"}).
		SetRaw([]byte{1}).
		SetPrivateFrom("testPrivateFrom").
		SetPrivateFor([]string{"testPrivateFor1", "testPrivateFor2"}).
		SetContextLabelsValue("testContextKey", "testContextValue")

	errs := b.Validate()
	assert.Len(t, errs, 5)
}

func TestEnvelope_InternalLabels(t *testing.T) {
	b := NewEnvelope().SetInternalLabelsValue("key", "value")
	assert.Equal(t, "value", b.GetInternalLabelsValue("key"), "Should be equal")
}

func TestEnvelope_From(t *testing.T) {
	from := ethcommon.HexToAddress("0x1")
	b := NewEnvelope().SetFrom(from)
	assert.Equal(t, from, *b.GetFrom(), "Should be equal")

	addr, err := b.GetFromAddress()
	assert.Equal(t, from, addr, "Should be equal")
	assert.NoError(t, err, "Should be nil")

	assert.Equal(t, from, b.MustGetFromAddress(), "Should be equal")
	assert.Equal(t, "0x0000000000000000000000000000000000000001", b.GetFromString(), "Should be equal")

	b.From = nil
	addr, err = b.GetFromAddress()
	assert.Equal(t, ethcommon.Address{}, addr, "Should be equal")
	assert.Error(t, err, "Should be not nil")
	assert.Equal(t, ethcommon.Address{}, b.MustGetFromAddress(), "Should be equal")
	assert.Equal(t, "", b.GetFromString(), "Should be equal")

	assert.Error(t, b.MustSetFromString("z"), "Should be not nil")
	assert.Error(t, b.SetFromString("z"), "Should be not nil")

	_ = b.MustSetFromString("0x1")
	assert.Equal(t, "0x0000000000000000000000000000000000000001", b.GetFromString(), "Should be equal")
}

func TestEnvelope_To(t *testing.T) {
	from := ethcommon.HexToAddress("0x1")
	b := NewEnvelope().SetTo(from)
	assert.Equal(t, from, *b.GetTo(), "Should be equal")

	addr, err := b.GetToAddress()
	assert.Equal(t, from, addr, "Should be equal")
	assert.NoError(t, err, "Should be nil")

	assert.Equal(t, from, b.MustGetToAddress(), "Should be equal")
	assert.Equal(t, "0x0000000000000000000000000000000000000001", b.GetToString(), "Should be equal")

	b.To = nil
	addr, err = b.GetToAddress()
	assert.Equal(t, ethcommon.Address{}, addr, "Should be equal")
	assert.Error(t, err, "Should be not nil")
	assert.Equal(t, ethcommon.Address{}, b.MustGetToAddress(), "Should be equal")
	assert.Equal(t, "", b.GetToString(), "Should be equal")

	assert.Error(t, b.MustSetToString("z"), "Should be not nil")
	assert.Error(t, b.SetToString("z"), "Should be not nil")

	_ = b.MustSetToString("0x1")
	assert.Equal(t, "0x0000000000000000000000000000000000000001", b.GetToString(), "Should be equal")
}

func TestEnvelope_Gas(t *testing.T) {
	b := NewEnvelope().SetGas(10)
	assert.Equal(t, uint64(10), *b.GetGas(), "Should be equal")
	assert.Equal(t, uint64(10), b.MustGetGasUint64(), "Should be equal")
	assert.Equal(t, "10", b.GetGasString(), "Should be equal")
	gas, err := b.GetGasUint64()
	assert.Equal(t, uint64(10), gas, "Should be equal")
	assert.NoError(t, err)

	b.Gas = nil
	gas, err = b.GetGasUint64()
	assert.Equal(t, uint64(0), gas, "Should be equal")
	assert.Error(t, err, "Should be not nil")
	assert.Equal(t, uint64(0), b.MustGetGasUint64(), "Should be equal")
	assert.Equal(t, "", b.GetGasString(), "Should be equal")

	err = b.SetGasString("12")
	assert.Equal(t, "12", b.GetGasString(), "Should be equal")
	assert.NoError(t, err)

	err = b.SetGasString("@")
	assert.Error(t, err)
}

func TestEnvelope_Nonce(t *testing.T) {
	b := NewEnvelope().SetNonce(10)
	assert.Equal(t, uint64(10), *b.GetNonce(), "Should be equal")
	assert.Equal(t, uint64(10), b.MustGetNonceUint64(), "Should be equal")
	assert.Equal(t, "10", b.GetNonceString(), "Should be equal")
	gas, err := b.GetNonceUint64()
	assert.Equal(t, uint64(10), gas, "Should be equal")
	assert.NoError(t, err)

	b.Nonce = nil
	gas, err = b.GetNonceUint64()
	assert.Equal(t, uint64(0), gas, "Should be equal")
	assert.Error(t, err, "Should be not nil")
	assert.Equal(t, uint64(0), b.MustGetNonceUint64(), "Should be equal")
	assert.Equal(t, "", b.GetNonceString(), "Should be equal")

	err = b.SetNonceString("12")
	assert.Equal(t, "12", b.GetNonceString(), "Should be equal")
	assert.NoError(t, err)

	err = b.SetNonceString("@")
	assert.Error(t, err)
}

func TestEnvelope_GasPrice(t *testing.T) {
	b := NewEnvelope().SetGasPrice(big.NewInt(10))
	assert.Equal(t, big.NewInt(10), b.GetGasPrice(), "Should be equal")
	assert.Equal(t, "10", b.GetGasPriceString(), "Should be equal")
	gasPrice, err := b.GetGasPriceBig()
	assert.Equal(t, big.NewInt(10), gasPrice, "Should be equal")
	assert.NoError(t, err)

	b.GasPrice = nil
	var nilGasPrice *big.Int
	gasPrice, err = b.GetGasPriceBig()
	assert.Equal(t, nilGasPrice, gasPrice, "Should be equal")
	assert.Error(t, err, "Should be not nil")
	assert.Equal(t, "", b.GetGasPriceString(), "Should be equal")

	err = b.SetGasPriceString("12")
	assert.Equal(t, "12", b.GetGasPriceString(), "Should be equal")
	assert.NoError(t, err)
	err = b.SetGasPriceString("@")
	assert.Error(t, err)
}

func TestEnvelope_Value(t *testing.T) {
	b := NewEnvelope().SetValue(big.NewInt(10))
	assert.Equal(t, big.NewInt(10), b.GetValue(), "Should be equal")
	assert.Equal(t, "10", b.GetValueString(), "Should be equal")
	value, err := b.GetValueBig()
	assert.Equal(t, big.NewInt(10), value, "Should be equal")
	assert.NoError(t, err)

	b.Value = nil
	var nilValue *big.Int
	value, err = b.GetValueBig()
	assert.Equal(t, nilValue, value, "Should be equal")
	assert.Error(t, err, "Should be not nil")
	assert.Equal(t, "", b.GetValueString(), "Should be equal")

	err = b.SetValueString("12")
	assert.Equal(t, "12", b.GetValueString(), "Should be equal")
	assert.NoError(t, err)
	err = b.SetValueString("@")
	assert.Error(t, err)
}

func TestEnvelope_Data(t *testing.T) {
	b := NewEnvelope()
	err := b.SetDataString("0x01")
	assert.NoError(t, err)
	assert.Equal(t, "0x01", b.GetData(), "Should be equal")
	assert.Equal(t, []byte{1}, b.MustGetDataBytes(), "Should be equal")

	b.Data = ""
	assert.Equal(t, []byte{}, b.MustGetDataBytes(), "Should be equal")

	_ = b.SetData([]byte{2})
	assert.Equal(t, []byte{2}, b.MustGetDataBytes(), "Should be equal")

	err = b.SetDataString("@")
	assert.Error(t, err)

	_ = b.MustSetDataString("@")
}

func TestEnvelope_Raw(t *testing.T) {
	b := NewEnvelope()
	err := b.SetRawString("0x01")
	assert.NoError(t, err)
	assert.Equal(t, "0x01", b.GetRaw(), "Should be equal")
	assert.Equal(t, "0x01", b.GetShortRaw(), "Should be equal")
	assert.Equal(t, []byte{1}, b.MustGetRawBytes(), "Should be equal")

	b.Raw = ""
	assert.Equal(t, []byte{}, b.MustGetRawBytes(), "Should be equal")

	_ = b.SetRaw([]byte{2})
	assert.Equal(t, []byte{2}, b.MustGetRawBytes(), "Should be equal")

	err = b.SetRawString("@")
	assert.Error(t, err)

	_ = b.MustSetRawString("@")
}

func TestEnvelope_TxHash(t *testing.T) {
	txHash := "0x9c1c1cd76f408145a2bd14ccfb16517ffb28dc99afc992503cff554c683828e3"
	b := NewEnvelope()
	err := b.SetTxHashString(txHash)
	assert.NoError(t, err)
	assert.Equal(t, txHash, b.GetTxHashString(), "Should be equal")
	assert.Equal(t, txHash, b.GetTxHash().Hex(), "Should be equal")
	assert.Equal(t, txHash, b.MustGetTxHashValue().Hex(), "Should be equal")

	hash, err := b.GetTxHashValue()
	assert.Equal(t, hash, b.MustGetTxHashValue(), "Should be equal")
	assert.NoError(t, err)

	b.TxHash = nil
	_, err = b.GetTxHashValue()
	assert.Error(t, err)
	assert.Empty(t, b.GetTxHashString())
	assert.Equal(t, ethcommon.Hash{}, b.MustGetTxHashValue(), "Should be equal")

	err = b.SetTxHashString("@")
	assert.Error(t, err)
}

func TestEnvelope_ChainID(t *testing.T) {
	b := NewEnvelope().SetChainID(big.NewInt(10))
	assert.Equal(t, big.NewInt(10), b.GetChainID(), "Should be equal")
	assert.Equal(t, "10", b.GetChainIDString(), "Should be equal")

	err := b.SetChainIDString("11")
	assert.NoError(t, err)

	_ = b.SetChainIDUint64(12)
	assert.Equal(t, "12", b.GetChainIDString(), "Should be equal")

	b.ChainID = nil
	assert.Empty(t, b.GetChainIDString())

	err = b.SetChainIDString("@")
	assert.Error(t, err)
}

func TestEnvelope_ChainName(t *testing.T) {
	b := NewEnvelope().SetChainName("test")
	assert.Equal(t, "test", b.GetChainName(), "Should be equal")
}

func TestEnvelope_ChainUUID(t *testing.T) {
	b := NewEnvelope().SetChainUUID("test")
	assert.Equal(t, "test", b.GetChainUUID(), "Should be equal")
}

func TestEnvelope_Contract(t *testing.T) {
	b := NewEnvelope().SetContractName("testContractName").SetContractTag("testContractTag").MustSetToString("0xto")
	assert.Equal(t, &abi.ContractId{Name: "testContractName", Tag: "testContractTag"}, b.GetContractID(), "Should be equal")
	assert.False(t, b.IsContractCreation())

	assert.Equal(t, "testContractName[testContractTag]", b.ShortContract(), "Should be equal")

	b.ContractTag = ""
	assert.Equal(t, "testContractName", b.ShortContract(), "Should be equal")

	b.ContractName = ""
	assert.Empty(t, b.ShortContract())
}

func TestEnvelope_MethodSignature(t *testing.T) {
	b := NewEnvelope().SetContractName("testContractName").SetContractTag("testContractTag").MustSetToString("0xto")
	assert.Equal(t, &abi.ContractId{Name: "testContractName", Tag: "testContractTag"}, b.GetContractID(), "Should be equal")
	assert.False(t, b.IsContractCreation())

	assert.Equal(t, "testContractName[testContractTag]", b.ShortContract(), "Should be equal")

	b.To = nil
	assert.True(t, b.IsContractCreation())

}

func TestEnvelope_Args(t *testing.T) {
	args := []string{"test", "test2"}
	b := NewEnvelope().SetArgs(args)

	assert.Equal(t, args, b.GetArgs(), "Should be equal")
}

func TestEnvelope_Private(t *testing.T) {
	args := []string{"test", "test2"}
	pFrom := "testFrom"
	b := NewEnvelope().SetPrivateFor(args).SetPrivateFrom(pFrom)
	assert.Equal(t, args, b.GetPrivateFor(), "Should be equal")
	assert.Equal(t, pFrom, b.GetPrivateFrom(), "Should be equal")
}

func TestEnvelope_TxRequest(t *testing.T) {
	b := NewEnvelope().
		SetID("jobUUID").
		SetHeadersValue("testHeaderKey", "testHeaderValue").
		SetChainName("chainName").
		MustSetFromString("0x1").
		MustSetToString("0x2").
		SetGas(11).
		SetGasPrice(big.NewInt(12)).
		SetValue(big.NewInt(13)).
		SetNonce(14).
		SetData([]byte{1}).
		SetContractName("testContractName").
		SetContractTag("testContractTag").
		SetMethodSignature("testMethodSignature").
		SetArgs([]string{"testArg1", "testArg2"}).
		SetRaw([]byte{2}).
		SetPrivateFrom("testPrivateFrom").
		SetPrivateFor([]string{"testPrivateFor1", "testPrivateFor2"}).
		SetContextLabelsValue("testContextKey", "testContextValue")

	req := &TxRequest{
		Id:      "jobUUID",
		Headers: map[string]string{"testHeaderKey": "testHeaderValue"},
		Chain:   "chainName",
		Params: &Params{
			From:            "0x0000000000000000000000000000000000000001",
			To:              "0x0000000000000000000000000000000000000002",
			Gas:             "11",
			GasPrice:        "12",
			Value:           "13",
			Nonce:           "14",
			Data:            "0x01",
			Contract:        "testContractName[testContractTag]",
			MethodSignature: "testMethodSignature",
			Args:            []string{"testArg1", "testArg2"},
			Raw:             "0x02",
			PrivateFor:      []string{"testPrivateFor1", "testPrivateFor2"},
			PrivateFrom:     "testPrivateFrom",
			PrivateTxType:   "",
			PrivacyGroupId:  "",
		},
		ContextLabels: map[string]string{"testContextKey": "testContextValue"},
	}

	assert.Equal(t, req, b.TxRequest(), "Should be equal")
}

func TestEnvelope_TxEnvelopeAsRequest(t *testing.T) {
	b := NewEnvelope().
		SetID("envelopeID").
		SetJobUUID("jobUUID").
		SetScheduleUUID("scheduleUUID").
		SetHeadersValue("testHeaderKey", "testHeaderValue").
		SetChainID(big.NewInt(1)).
		SetChainName("chainName").
		SetChainUUID("testChainUUID").
		MustSetFromString("0x1").
		MustSetToString("0x2").
		SetGas(11).
		SetGasPrice(big.NewInt(12)).
		SetValue(big.NewInt(13)).
		SetNonce(14).
		SetData([]byte{1}).
		SetContractName("testContractName").
		SetContractTag("testContractTag").
		SetMethodSignature("testMethodSignature").
		SetArgs([]string{"testArg1", "testArg2"}).
		SetRaw([]byte{2}).
		SetPrivateFrom("testPrivateFrom").
		SetPrivateFor([]string{"testPrivateFor1", "testPrivateFor2"}).
		SetContextLabelsValue("testContextKey", "testContextValue").
		MustSetTxHashString("0x2d6a7b0f6adeff38423d4c62cd8b6ccb708ddad85da5d3d06756ad4d8a04a6a2").
		AppendError(errors.DataError("testError"))

	txEnvelopeReq := &TxEnvelope{
		Msg: &TxEnvelope_TxRequest{
			TxRequest: &TxRequest{
				Id:      "envelopeID",
				Headers: map[string]string{"testHeaderKey": "testHeaderValue"},
				Chain:   "chainName",
				Params: &Params{
					From:            "0x0000000000000000000000000000000000000001",
					To:              "0x0000000000000000000000000000000000000002",
					Gas:             "11",
					GasPrice:        "12",
					Value:           "13",
					Nonce:           "14",
					Data:            "0x01",
					Contract:        "testContractName[testContractTag]",
					MethodSignature: "testMethodSignature",
					Args:            []string{"testArg1", "testArg2"},
					Raw:             "0x02",
					PrivateFor:      []string{"testPrivateFor1", "testPrivateFor2"},
					PrivateFrom:     "testPrivateFrom",
					PrivateTxType:   "",
					PrivacyGroupId:  "",
				},
				ContextLabels: map[string]string{"testContextKey": "testContextValue"},
			},
		},
		InternalLabels: map[string]string{
			ChainIDLabel:      "1",
			TxHashLabel:       "0x2d6a7b0f6adeff38423d4c62cd8b6ccb708ddad85da5d3d06756ad4d8a04a6a2",
			ChainUUIDLabel:    "testChainUUID",
			ScheduleUUIDLabel: "scheduleUUID",
			JobUUIDLabel:      "jobUUID",
		},
	}

	txEnvelopeRes := &TxEnvelope{
		Msg: &TxEnvelope_TxResponse{
			TxResponse: &TxResponse{
				Id:            "envelopeID",
				JobUUID:       "jobUUID",
				Headers:       map[string]string{"testHeaderKey": "testHeaderValue"},
				ContextLabels: map[string]string{"testContextKey": "testContextValue"},
				Transaction: &ethereum.Transaction{
					From:     "0x0000000000000000000000000000000000000001",
					Nonce:    "14",
					To:       "0x0000000000000000000000000000000000000002",
					Value:    "13",
					Gas:      "11",
					GasPrice: "12",
					Data:     "0x01",
					Raw:      "0x02",
					TxHash:   "0x2d6a7b0f6adeff38423d4c62cd8b6ccb708ddad85da5d3d06756ad4d8a04a6a2",
				},
				Chain: "chainName",
				Errors: []*ierror.Error{
					{Message: "testError", Code: 270336},
				},
			},
		},
		InternalLabels: map[string]string{
			ChainIDLabel:      "1",
			TxHashLabel:       "0x2d6a7b0f6adeff38423d4c62cd8b6ccb708ddad85da5d3d06756ad4d8a04a6a2",
			ChainUUIDLabel:    "testChainUUID",
			ScheduleUUIDLabel: "scheduleUUID",
			JobUUIDLabel:      "jobUUID",
		},
	}

	assert.Equal(t, txEnvelopeReq, b.TxEnvelopeAsRequest(), "Should be equal")
	assert.Equal(t, txEnvelopeRes, b.TxEnvelopeAsResponse(), "Should be equal")
}

func TestEnvelope_TxResponse(t *testing.T) {
	b := NewEnvelope().
		SetID("envelopeID").
		SetJobUUID("jobUUID").
		SetScheduleUUID("scheduleUUID").
		SetHeadersValue("testHeaderKey", "testHeaderValue").
		SetChainName("chainName").
		MustSetFromString("0x1").
		MustSetToString("0x2").
		SetGas(11).
		SetGasPrice(big.NewInt(12)).
		SetValue(big.NewInt(13)).
		SetNonce(14).
		SetData([]byte{1}).
		SetContractName("testContractName").
		SetContractTag("testContractTag").
		SetMethodSignature("testMethodSignature").
		SetArgs([]string{"testArg1", "testArg2"}).
		SetRaw([]byte{2}).
		MustSetTxHashString("0x2d6a7b0f6adeff38423d4c62cd8b6ccb708ddad85da5d3d06756ad4d8a04a6a2").
		SetPrivateFrom("testPrivateFrom").
		SetPrivateFor([]string{"testPrivateFor1", "testPrivateFor2"}).
		SetContextLabelsValue("testContextKey", "testContextValue").
		AppendError(errors.DataError("testError"))

	res := &TxResponse{
		Headers:       map[string]string{"testHeaderKey": "testHeaderValue"},
		Id:            "envelopeID",
		JobUUID:       "jobUUID",
		ContextLabels: map[string]string{"testContextKey": "testContextValue"},
		Transaction: &ethereum.Transaction{
			From:     "0x0000000000000000000000000000000000000001",
			Nonce:    "14",
			To:       "0x0000000000000000000000000000000000000002",
			Value:    "13",
			Gas:      "11",
			GasPrice: "12",
			Data:     "0x01",
			Raw:      "0x02",
			TxHash:   "0x2d6a7b0f6adeff38423d4c62cd8b6ccb708ddad85da5d3d06756ad4d8a04a6a2",
		},
		Chain:   "chainName",
		Receipt: b.Receipt,
		Errors:  b.Errors,
	}

	assert.Equal(t, res, b.TxResponse(), "Should be equal")
}

func TestPrivateValidation(t *testing.T) {
	b := NewEnvelope().
		SetID("9f8708ad-8019-4533-9690-6495cc79a03c").
		SetJobType(JobType_ETH_ORION_EEA_TX).
		SetPrivateFor([]string{"kAbelwaVW7okoEn1+okO+AbA4Hhz/7DaCOWVQz9nx5M="}).
		SetPrivacyGroupID("kAbelwaVW7okoEn1+okO+AbA4Hhz/7DaCOWVQz9nx5M=")

	assert.Len(t, b.Validate(), 1)
	if len(b.Validate()) == 1 {
		assert.Equal(t, errors.DataError("privacyGroupId and privateFor fields are mutually exclusive"), b.Validate()[0])
	}

	b2 := NewEnvelope().
		SetID("9f8708ad-8019-4533-9690-6495cc79a03c").
		SetJobType(JobType_ETH_TESSERA_PRIVATE_TX).
		SetPrivacyGroupID("kAbelwaVW7okoEn1+okO+AbA4Hhz/7DaCOWVQz9nx5M=")

	assert.Len(t, b2.Validate(), 0)

	b3 := NewEnvelope().
		SetID("9f8708ad-8019-4533-9690-6495cc79a03c").
		SetJobType(JobType_ETH_TESSERA_PRIVATE_TX).
		SetPrivateFor([]string{"kAbelwaVW7okoEn1+okO+AbA4Hhz/7DaCOWVQz9nx5M="})

	assert.Len(t, b3.Validate(), 0)
}

func TestPartitionKey(t *testing.T) {
	b := NewEnvelope().
		SetChainID(big.NewInt(10)).
		SetChainName("testChain").
		MustSetFromString("0x1").
		SetID("9f8708ad-8019-4533-9690-6495cc79a03c").
		SetPrivacyGroupID("kAbelwaVW7okoEn1+okO+AbA4Hhz/7DaCOWVQz9nx5M=")
	assert.Equal(t, "0x0000000000000000000000000000000000000001@10", b.PartitionKey())

	b2 := NewEnvelope().
		SetChainID(big.NewInt(11)).
		SetChainName("testChain").
		MustSetFromString("0x1").
		SetID("9f8708ad-8019-4533-9690-6495cc79a03c").
		SetJobType(JobType_ETH_ORION_EEA_TX).
		SetPrivacyGroupID("kAbelwaVW7okoEn1+okO+AbA4Hhz/7DaCOWVQz9nx5M=")
	assert.Equal(t, "0x0000000000000000000000000000000000000001@orion-kAbelwaVW7okoEn1+okO+AbA4Hhz/7DaCOWVQz9nx5M=@11", b2.PartitionKey())

	b3 := NewEnvelope().
		SetChainID(big.NewInt(12)).
		SetChainName("testChain").
		MustSetFromString("0x1").
		SetID("9f8708ad-8019-4533-9690-6495cc79a03c").
		SetJobType(JobType_ETH_ORION_EEA_TX).
		SetPrivateFor([]string{"kAbelwaVW7okoEn1+okO+AbA4Hhz/7DaCOWVQz9nx5M="})
	assert.Equal(t, "0x0000000000000000000000000000000000000001@orion-a3ce4ff3ac5af3264fd8ae06af53ed9e@12", b3.PartitionKey())
}
