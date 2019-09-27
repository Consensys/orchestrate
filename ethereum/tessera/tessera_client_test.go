package tessera

import (
	"encoding/base64"
	"errors"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/mocks"
)

const chainID = 888

var chainIDBigInt = big.NewInt(int64(chainID))

var errTest = errors.New("test error")

var rawTx = []byte{1, 2, 3}
var privateFrom = "0x01"
var expectedStorerawRequest = map[string]string{
	"payload": "AQID",
	"from":    privateFrom,
}

var resultTxHash = []byte{1}
var resultTxHashBase64 = base64.StdEncoding.EncodeToString(resultTxHash)
var resultStoreRawResponse = createSendRawResponse(resultTxHashBase64)

var enclaveClient *EnclaveClient
var ctrl *gomock.Controller
var mockEnclaveEndpoint *mocks.MockEnclaveEndpoint

func setupTest(t *testing.T) {
	ctrl = gomock.NewController(t)
	enclaveClient = NewEnclaveClient()
	mockEnclaveEndpoint = mocks.NewMockEnclaveEndpoint(ctrl)
}

func TestStoreRawTransaction(t *testing.T) {
	setupTest(t)
	defer ctrl.Finish()

	setMockClient(mockEnclaveEndpoint)
	mockSuccessfulStoreRawResult(resultStoreRawResponse)

	hash, err := enclaveClient.StoreRaw(chainIDBigInt.String(), rawTx, privateFrom)
	assert.NoError(t, err)
	assert.Equal(t, resultTxHash, hash)
}

func TestReturnErrorIfCannotGetEndpointForChain(t *testing.T) {
	setupTest(t)
	defer ctrl.Finish()

	hash, err := enclaveClient.StoreRaw(chainIDBigInt.String(), rawTx, privateFrom)

	assertError(t, hash, err, "no Tessera endpoint for chain id: 888")
}

func TestQuorumRawPrivateTransactionWhenRPCCallFails(t *testing.T) {
	setupTest(t)
	defer ctrl.Finish()

	setMockClient(mockEnclaveEndpoint)
	mockUnsuccessfulStoreRawResult(errTest)

	hash, err := enclaveClient.StoreRaw(chainIDBigInt.String(), rawTx, privateFrom)

	assertError(t, hash, err, "failed to send a request to Tessera enclave: test error")
}

func TestReturnInvalidBase64(t *testing.T) {
	setupTest(t)
	defer ctrl.Finish()

	setMockClient(mockEnclaveEndpoint)

	invalidBase64 := "a"
	mockSuccessfulStoreRawResult(createSendRawResponse(invalidBase64))

	hash, err := enclaveClient.StoreRaw(chainIDBigInt.String(), rawTx, privateFrom)

	assertError(t, hash, err, "failed to decode base64 encoded string in the 'storeraw' response: illegal base64 data at input byte 0")
}

func TestGetStatus(t *testing.T) {
	setupTest(t)
	defer ctrl.Finish()

	setMockClient(mockEnclaveEndpoint)
	mockUpcheckResult("status", nil)

	status, err := enclaveClient.GetStatus(chainIDBigInt.String())
	assert.NoError(t, err)
	assert.Equal(t, "status", status)
}

func createSendRawResponse(base64Hash string) StoreRawResponse {
	return StoreRawResponse{
		Key: base64Hash,
	}
}

func setMockClient(tesseraEndpoint EnclaveEndpoint) {
	enclaveClient.AddClient(chainIDBigInt.String(), tesseraEndpoint)
}

func assertError(t *testing.T, hash []byte, err error, errMessage string) {
	assert.EqualError(t, err, errMessage)
	assert.Nil(t, nil, hash)
}

func mockSuccessfulStoreRawResult(response interface{}) {
	mockStoreRawResult(response, nil)
}

func mockUnsuccessfulStoreRawResult(err error) {
	mockStoreRawResult(StoreRawResponse{}, err)
}

func mockStoreRawResult(response interface{}, err error) {
	mockEnclaveEndpoint.
		EXPECT().
		PostRequest("storeraw", expectedStorerawRequest, gomock.Any()).
		Return(err).
		SetArg(2, response).
		Times(1)
}

func mockUpcheckResult(response string, err error) {
	mockEnclaveEndpoint.
		EXPECT().
		GetRequest("upcheck").
		Return(response, err).
		Times(1)
}
