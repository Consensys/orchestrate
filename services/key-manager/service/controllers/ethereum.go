package controllers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/key-manager/use-cases"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"

	"github.com/gorilla/mux"
	jsonutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/httputil"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/service/formatters"
)

const ethAccountPath = "/ethereum/accounts"

type EthereumController struct {
	vault    store.Vault
	useCases usecases.ETHUseCases
}

func NewEthereumController(vault store.Vault, useCases usecases.ETHUseCases) *EthereumController {
	return &EthereumController{vault: vault, useCases: useCases}
}

// Add routes to router
func (c *EthereumController) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/ethereum/namespaces").HandlerFunc(c.listNamespaces)
	router.Methods(http.MethodPost).Path(ethAccountPath).HandlerFunc(c.createAccount)
	router.Methods(http.MethodGet).Path(ethAccountPath).HandlerFunc(c.listAccounts)
	router.Methods(http.MethodPost).Path(ethAccountPath + "/import").HandlerFunc(c.importAccount)
	router.Methods(http.MethodGet).Path(ethAccountPath + "/{address}").HandlerFunc(c.getAccount)
	router.Methods(http.MethodPost).Path(ethAccountPath + "/{address}/sign").HandlerFunc(c.signPayload)
	router.Methods(http.MethodPost).Path(ethAccountPath + "/{address}/sign-transaction").HandlerFunc(c.signTransaction)
	router.Methods(http.MethodPost).Path(ethAccountPath + "/{address}/sign-eea-transaction").HandlerFunc(c.signEEA)
	router.Methods(http.MethodPost).Path(ethAccountPath + "/{address}/sign-quorum-private-transaction").HandlerFunc(c.signQuorumPrivate)
	router.Methods(http.MethodPost).Path(ethAccountPath + "/{address}/sign-typed-data").HandlerFunc(c.signTypedData)
	router.Methods(http.MethodPost).Path(ethAccountPath + "/verify-signature").HandlerFunc(c.verifySignature)
	router.Methods(http.MethodPost).Path(ethAccountPath + "/verify-typed-data-signature").HandlerFunc(c.verifyTypedDataSignature)
}

// @Summary Creates a new Ethereum Account
// @Description Creates a new private key, stores it in the Vault and generates a public key given a chosen elliptic curve
// @Accept json
// @Produce json
// @Param request body ethereum.CreateETHAccountRequest true "Ethereum account creation request"
// @Success 200 {object} ethereum.ETHAccountResponse "Created Ethereum account"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /ethereum/accounts [post]
func (c *EthereumController) createAccount(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	ethAccountRequest := &types.CreateETHAccountRequest{}
	err := jsonutils.UnmarshalBody(request.Body, ethAccountRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	accountResponse, err := c.vault.ETHCreateAccount(ethAccountRequest.Namespace)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatETHAccountResponse(accountResponse))
}

// @Summary List Ethereum Accounts
// @Description List stored ethereum account in the Vault
// @Produce json
// @Param namespace query string false "namespace where key is stored"
// @Success 200 {object} []string "List of ethereum public accounts"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /ethereum/accounts [get]
func (c *EthereumController) listAccounts(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	namespace := req.URL.Query().Get("namespace")

	accountAddrs, err := c.vault.ETHListAccounts(namespace)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(accountAddrs)
}

// @Summary List Ethereum Namespaces
// @Description List ethereum namespaces in the Vault
// @Produce json
// @Success 200 {object} []string "List of ethereum public namespaces"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /ethereum/namespaces [get]
func (c *EthereumController) listNamespaces(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	namespaces, err := c.vault.ETHListNamespaces()
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(namespaces)
}

// @Summary Fetch Ethereum Account
// @Description Get selected stored ethereum account in the Vault
// @Produce json
// @Success 200 {object} ethereum.ETHAccountResponse "Ethereum account"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /ethereum/accounts/{address} [get]
func (c *EthereumController) getAccount(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	namespace := req.URL.Query().Get("namespace")
	address, err := utils.ParseHexToMixedCaseEthAddress(mux.Vars(req)["publicKey"])
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	ethAcc, err := c.vault.ETHGetAccount(address, namespace)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}
	if ethAcc == nil {
		httputil.WriteHTTPErrorResponse(rw, errors.NotFoundError("account not found"))
		return
	}

	_ = json.NewEncoder(rw).Encode(ethAcc)
}

// @Summary Imports an Ethereum Account
// @Description Imports a private key, stores it in the Vault and generates a public key given a chosen elliptic curve
// @Accept json
// @Produce json
// @Param request body ethereum.ImportETHAccountRequest true "Ethereum account import request"
// @Success 200 {object} ethereum.ETHAccountResponse "Imported Ethereum account"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 422 {object} httputil.ErrorResponse "Invalid private key"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /ethereum/accounts/import [post]
func (c *EthereumController) importAccount(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	ethAccountRequest := &types.ImportETHAccountRequest{}
	err := jsonutils.UnmarshalBody(request.Body, ethAccountRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	accountResponse, err := c.vault.ETHImportAccount(ethAccountRequest.Namespace, ethAccountRequest.PrivateKey)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatETHAccountResponse(accountResponse))
}

// @Summary Signs an arbitrary message using an existing Ethereum account
// @Description Signs an arbitrary message using ECDSA and the private key of an existing Ethereum account
// @Accept json
// @Produce text/plain
// @Param request body keymanager.PayloadRequest true "Payload to sign"
// @Success 200 {string} string "Signed payload"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 404 {object} httputil.ErrorResponse "Account not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /ethereum/accounts/{address}/sign [post]
func (c *EthereumController) signPayload(rw http.ResponseWriter, request *http.Request) {
	signRequest := &keymanager.SignPayloadRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	address, err := utils.ParseHexToMixedCaseEthAddress(mux.Vars(request)["address"])
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	signature, err := c.vault.ETHSign(address, signRequest.Namespace, signRequest.Data)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(signature))
}

// @Summary Signs an Ethereum transaction using an existing account
// @Description Signs an Ethereum transaction using ECDSA and the private key of an existing account
// @Accept json
// @Produce text/plain
// @Param request body ethereum.SignETHTransactionRequest true "Ethereum transaction to sign"
// @Success 200 {string} string "Signed payload"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 404 {object} httputil.ErrorResponse "Account not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /ethereum/accounts/{address}/sign-transaction [post]
func (c *EthereumController) signTransaction(rw http.ResponseWriter, request *http.Request) {
	signRequest := &types.SignETHTransactionRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	address, err := utils.ParseHexToMixedCaseEthAddress(mux.Vars(request)["address"])
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	signature, err := c.vault.ETHSignTransaction(address, signRequest)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(signature))
}

// @Summary Signs a Quorum private transaction using an existing account
// @Description Signs a Quorum private transaction using ECDSA and the private key of an existing account
// @Accept json
// @Produce text/plain
// @Param request body ethereum.SignQuorumPrivateTransactionRequest true "Quorum private transaction to sign"
// @Success 200 {string} string "Signed payload"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 404 {object} httputil.ErrorResponse "Account not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /ethereum/accounts/{address}/sign-quorum-private-transaction [post]
func (c *EthereumController) signQuorumPrivate(rw http.ResponseWriter, request *http.Request) {
	signRequest := &types.SignQuorumPrivateTransactionRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	address, err := utils.ParseHexToMixedCaseEthAddress(mux.Vars(request)["address"])
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	signature, err := c.vault.ETHSignQuorumPrivateTransaction(address, signRequest)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(signature))
}

// @Summary Signs an EEA private transaction using an existing account
// @Description Signs an EEA private transaction using ECDSA and the private key of an existing account
// @Accept json
// @Produce text/plain
// @Param request body ethereum.SignQuorumPrivateTransactionRequest true "EEA private transaction to sign"
// @Success 200 {string} string "Signed payload"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 404 {object} httputil.ErrorResponse "Account not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /ethereum/accounts/{address}/sign-eea-transaction [post]
func (c *EthereumController) signEEA(rw http.ResponseWriter, request *http.Request) {
	signRequest := &types.SignEEATransactionRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = signRequest.Validate()
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	address, err := utils.ParseHexToMixedCaseEthAddress(mux.Vars(request)["address"])
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	signature, err := c.vault.ETHSignEEATransaction(address, signRequest)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(signature))
}

// @Summary Signs typed data using an existing account following the EIP-712 standard
// @Description Signs typed data using ECDSA and the private key of an existing account following the EIP-712 standard
// @Accept json
// @Produce text/plain
// @Param request body ethereum.SignTypedDataRequest{domainSeparator=types.DomainSeparator} true "Typed data to sign"
// @Success 200 {string} string "Signed payload"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 404 {object} httputil.ErrorResponse "Account not found"
// @Failure 422 {object} httputil.ErrorResponse "Invalid parameters"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /ethereum/accounts/{address}/sign-typed-data [post]
func (c *EthereumController) signTypedData(rw http.ResponseWriter, request *http.Request) {
	signRequest := &types.SignTypedDataRequest{}
	err := jsonutils.UnmarshalBody(request.Body, signRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	address, err := utils.ParseHexToMixedCaseEthAddress(mux.Vars(request)["address"])
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	typedData := formatters.FormatSignTypedDataRequest(signRequest)
	signature, err := c.useCases.SignTypedData().Execute(request.Context(), address, signRequest.Namespace, typedData)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(signature))
}

// @Summary Verifies the signature of a typed data message following the EIP-712 standard
// @Description Verifies if a typed data message has been signed by the Ethereum account passed as argument following the EIP-712 standard
// @Accept json
// @Param request body ethereum.SignTypedDataRequest{domainSeparator=types.DomainSeparator} true "Typed data to sign"
// @Success 204
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 404 {object} httputil.ErrorResponse "Account not found"
// @Failure 422 {object} httputil.ErrorResponse "Invalid parameters"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /ethereum/accounts/verify-typed-data-signature [post]
func (c *EthereumController) verifyTypedDataSignature(rw http.ResponseWriter, request *http.Request) {
	verifyRequest := &types.VerifyTypedDataRequest{}
	err := jsonutils.UnmarshalBody(request.Body, verifyRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	address, err := utils.ParseHexToMixedCaseEthAddress(verifyRequest.Address)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	typedData := formatters.FormatSignTypedDataRequest(&verifyRequest.TypedData)
	err = c.useCases.VerifyTypedDataSignature().Execute(request.Context(), address, verifyRequest.Signature, typedData)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

// @Summary Verifies the signature of a message
// @Description Verifies if a message has been signed by the Ethereum account passed as argument
// @Accept json
// @Param request body ethereum.VerifySignatureRequest true "signature and message to verify"
// @Success 204
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 422 {object} httputil.ErrorResponse "Failed to verify"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /ethereum/accounts/verify-signature [post]
func (c *EthereumController) verifySignature(rw http.ResponseWriter, request *http.Request) {
	verifyRequest := &types.VerifyPayloadRequest{}
	err := jsonutils.UnmarshalBody(request.Body, verifyRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	address, err := utils.ParseHexToMixedCaseEthAddress(verifyRequest.Address)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.useCases.VerifySignature().Execute(request.Context(), address, verifyRequest.Signature, verifyRequest.Data)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
