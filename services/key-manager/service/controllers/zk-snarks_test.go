// +build unit

package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/key-manager/use-cases"
	mocks2 "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/key-manager/use-cases/mocks"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/service/formatters"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store/mocks"
)

type zksCtrlTestSuite struct {
	suite.Suite
	vault             *mocks.MockVault
	verifySignatureUC *mocks2.MockVerifyZKSSignatureUseCase
	router            *mux.Router
}

var _ usecases.ZKSUseCases = &zksCtrlTestSuite{}

const (
	testPublicAddress = "16551006344732991963827342392501535507890487822471009342749102663105305595515"
)

func (s zksCtrlTestSuite) VerifySignature() usecases.VerifyZKSSignatureUseCase {
	return s.verifySignatureUC
}

func TestZKSController(t *testing.T) {
	s := new(zksCtrlTestSuite)
	suite.Run(t, s)
}

func (s *zksCtrlTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	s.verifySignatureUC = mocks2.NewMockVerifyZKSSignatureUseCase(ctrl)
	s.vault = mocks.NewMockVault(ctrl)
	s.router = mux.NewRouter()

	controller := NewZKSController(s.vault, s)
	controller.Append(s.router)
}

func (s *zksCtrlTestSuite) TestZKSController_Create() {
	s.T().Run("should execute request successfully", func(t *testing.T) {
		createAccountRequest := testutils.FakeCreateZKSAccountRequest()
		requestBytes, _ := json.Marshal(createAccountRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, zksAccountPath, bytes.NewReader(requestBytes))

		fakeAccount := testutils.FakeZKSAccount()

		s.vault.EXPECT().ZKSCreateAccount(createAccountRequest.Namespace).Return(fakeAccount, nil)

		s.router.ServeHTTP(rw, httpRequest)

		response := formatters.FormatZKSAccountResponse(fakeAccount)
		expectedBody, _ := json.Marshal(response)
		assert.Equal(t, string(expectedBody)+"\n", rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should fail with correct error code if use case fails", func(t *testing.T) {
		createAccountRequest := testutils.FakeCreateZKSAccountRequest()
		requestBytes, _ := json.Marshal(createAccountRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, zksAccountPath, bytes.NewReader(requestBytes))

		s.vault.EXPECT().ZKSCreateAccount(gomock.Any()).Return(nil, errors.HashicorpVaultConnectionError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusInternalServerError, rw.Code)
	})
}

func (s *zksCtrlTestSuite) TestZKSController_Sign() {
	url := fmt.Sprintf("%s/%v/sign", zksAccountPath, testPublicAddress)

	s.T().Run("should execute request successfully", func(t *testing.T) {
		signature := "signature"
		payloadRequest := &keymanager.SignPayloadRequest{
			Data:      "my data to sign",
			Namespace: "namespace",
		}
		requestBytes, _ := json.Marshal(payloadRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(requestBytes))

		s.vault.EXPECT().
			ZKSSign(testPublicAddress, payloadRequest.Namespace, payloadRequest.Data).
			Return(signature, nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(t, signature, rw.Body.String())
		assert.Equal(t, http.StatusOK, rw.Code)
	})

	s.T().Run("should fail with correct error code if use case fails", func(t *testing.T) {
		payloadRequest := &keymanager.SignPayloadRequest{
			Data:      "my data to sign",
			Namespace: "namespace",
		}
		requestBytes, _ := json.Marshal(payloadRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(requestBytes))

		s.vault.EXPECT().
			ZKSSign(testPublicAddress, payloadRequest.Namespace, payloadRequest.Data).
			Return("", errors.InvalidParameterError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusUnprocessableEntity, rw.Code)
	})
}

func (s *zksCtrlTestSuite) TestZKSController_VerifySignature() {
	url := zksAccountPath + "/verify-signature"
	
	s.T().Run("should execute request successfully", func(t *testing.T) {
		verifyRequest := testutils.FakeZKSVerifyPayloadRequest()
		requestBytes, _ := json.Marshal(verifyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(requestBytes))

		s.verifySignatureUC.EXPECT().
			Execute(gomock.Any(), verifyRequest.PublicKey, verifyRequest.Signature, verifyRequest.Data).
			Return(nil)

		s.router.ServeHTTP(rw, httpRequest)

		assert.Equal(t, http.StatusNoContent, rw.Code)
	})
	
	s.T().Run("should fail with correct error code if use case fails", func(t *testing.T) {
		verifyRequest := testutils.FakeZKSVerifyPayloadRequest()
		requestBytes, _ := json.Marshal(verifyRequest)

		rw := httptest.NewRecorder()
		httpRequest := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(requestBytes))

		s.verifySignatureUC.EXPECT().
			Execute(gomock.Any(), verifyRequest.PublicKey, verifyRequest.Signature, verifyRequest.Data).
			Return(errors.InvalidParameterError("error"))

		s.router.ServeHTTP(rw, httpRequest)
		assert.Equal(t, http.StatusUnprocessableEntity, rw.Code)
	})
}
