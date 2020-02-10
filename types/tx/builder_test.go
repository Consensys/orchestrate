package tx

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	ierror "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/error"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/types/ethereum"
	"reflect"
	"testing"
)

func TestNewBuilder(t *testing.T) {
	b := NewBuilder()
	assert.NotNil(t, b, "Should not be nil")
	assert.NotNil(t, b.GetHeaders(), "Should not be nil")
	assert.NotNil(t, b.GetContextLabels(), "Should not be nil")
	assert.NotNil(t, b.GetErrors(), "Should not be nil")
	assert.NotNil(t, b.GetInternalLabels(), "Should not be nil")
}

func TestBuilder_SetID(t *testing.T) {
	id := uuid.NewV4().String()
	b := NewBuilder().SetID(id)
	assert.Equal(t, id, b.GetID(), "Should be equal")
}

func TestBuilder_GetErrors(t *testing.T) {
	testError := errors.FromError(fmt.Errorf("test"))
	b := NewBuilder().AppendError(testError)
	assert.True(t, reflect.DeepEqual(b.GetErrors(), []*ierror.Error{testError}), "Should be equal")
}

func TestBuilder_Error(t *testing.T) {
	b := NewBuilder()
	assert.Empty(t, b.Error(), "Should be equal")

	testError := errors.FromError(fmt.Errorf("test"))
	b = NewBuilder().AppendError(testError)
	assert.Equal(t, "[\"FF000@: test\"]", b.Error(), "Should be equal")
}

func TestBuilder_AppendError(t *testing.T) {
	testError := errors.FromError(fmt.Errorf("test"))
	b := NewBuilder().AppendError(testError)
	assert.True(t, reflect.DeepEqual(b.Errors, []*ierror.Error{testError}), "Should be equal")
}

func TestBuilder_AppendErrors(t *testing.T) {
	testErrors := []*ierror.Error{errors.FromError(fmt.Errorf("test1")), errors.FromError(fmt.Errorf("test2"))}
	b := NewBuilder().AppendErrors(testErrors)
	assert.True(t, reflect.DeepEqual(b.Errors, testErrors), "Should be equal")
}

func TestBuilder_SetReceipt(t *testing.T) {
	receipt := &ethereum.Receipt{TxHash: "test"}
	b := NewBuilder().SetReceipt(receipt)
	assert.True(t, reflect.DeepEqual(b.GetReceipt(), receipt), "Should be equal")

}

func TestBuilder_GetMethod(t *testing.T) {
	b := NewBuilder()
	assert.Equal(t, b.GetMethod(), Method_ETH_SENDRAWTRANSACTION, "Should be equal")
}

func TestBuilder_SetMethod(t *testing.T) {
	b := NewBuilder().SetMethod(Method_EEA_SENDPRIVATETRANSACTION)
	assert.Equal(t, b.GetMethod(), Method_EEA_SENDPRIVATETRANSACTION, "Should be equal")
}

func TestBuilder_IsMethod(t *testing.T) {
	b := NewBuilder()
	assert.True(t, b.IsEthSendRawTransaction(), "Should be equal")
	assert.False(t, b.IsEthSendPrivateTransaction(), "Should be equal")
	assert.False(t, b.IsEthSendRawPrivateTransaction(), "Should be equal")
	assert.False(t, b.IsEeaSendPrivateTransaction(), "Should be equal")

	_ = b.SetMethod(Method_ETH_SENDPRIVATETRANSACTION)
	assert.False(t, b.IsEthSendRawTransaction(), "Should be equal")
	assert.True(t, b.IsEthSendPrivateTransaction(), "Should be equal")
	assert.False(t, b.IsEthSendRawPrivateTransaction(), "Should be equal")
	assert.False(t, b.IsEeaSendPrivateTransaction(), "Should be equal")

	_ = b.SetMethod(Method_ETH_SENDRAWPRIVATETRANSACTION)
	assert.False(t, b.IsEthSendRawTransaction(), "Should be equal")
	assert.False(t, b.IsEthSendPrivateTransaction(), "Should be equal")
	assert.True(t, b.IsEthSendRawPrivateTransaction(), "Should be equal")
	assert.False(t, b.IsEeaSendPrivateTransaction(), "Should be equal")

	_ = b.SetMethod(Method_EEA_SENDPRIVATETRANSACTION)
	assert.False(t, b.IsEthSendRawTransaction(), "Should be equal")
	assert.False(t, b.IsEthSendPrivateTransaction(), "Should be equal")
	assert.False(t, b.IsEthSendRawPrivateTransaction(), "Should be equal")
	assert.True(t, b.IsEeaSendPrivateTransaction(), "Should be equal")
}

func TestBuilder_Carrier(t *testing.T) {
	b := NewBuilder()
}

func TestBuilder_OnlyWarnings(t *testing.T) {
	testError := errors.FromError(fmt.Errorf("test"))
	b := NewBuilder().AppendError(testError)
	assert.False(t, b.OnlyWarnings(), "Should be equal")

	testWarning := errors.Warningf("test")
	b = NewBuilder().AppendError(testWarning)
	assert.True(t, b.OnlyWarnings(), "Should be equal")
}

func TestBuilder_Headers(t *testing.T) {
	b := NewBuilder().SetHeadersValue("key", "value")
	assert.Equal(t, "value", b.GetHeadersValue("key"), "Should be equal")
}

func TestBuilder_ContextLabels(t *testing.T) {
	b := NewBuilder().SetContextLabelsValue("key", "value")
	assert.Equal(t, "value", b.GetContextLabelsValue("key"), "Should be equal")
}

func TestBuilder_InternalLabels(t *testing.T) {
	b := NewBuilder().SetInternalLabelsValue("key", "value")
	assert.Equal(t, "value", b.GetInternalLabelsValue("key"), "Should be equal")
}
