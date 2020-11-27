package controllers

import (
	"encoding/json"
	"net/http"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/utils"

	"github.com/gorilla/mux"
	jsonutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/httputil"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager/ethereum"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/service/formatters"
)

const Path = "/ethereum/accounts"

type EthereumController struct {
	vault store.Vault
}

func NewEthereumController(vault store.Vault) *EthereumController {
	return &EthereumController{vault: vault}
}

// Add routes to router
func (c *EthereumController) Append(router *mux.Router) {
	router.Methods(http.MethodPost).Path(Path).HandlerFunc(c.createAccount)
	router.Methods(http.MethodPost).Path(Path + "/import").HandlerFunc(c.importAccount)
	router.Methods(http.MethodPost).Path(Path + "/{address}/sign").HandlerFunc(c.signPayload)
	router.Methods(http.MethodPost).Path(Path + "/{address}/sign-transaction").HandlerFunc(c.signTransaction)
	router.Methods(http.MethodPost).Path(Path + "/{address}/sign-eea-transaction").HandlerFunc(c.signEEA)
	router.Methods(http.MethodPost).Path(Path + "/{address}/sign-quorum-private-transaction").HandlerFunc(c.signQuorumPrivate)
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
	signRequest := &keymanager.PayloadRequest{}
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
