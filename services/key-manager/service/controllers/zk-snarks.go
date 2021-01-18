package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	jsonutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/encoding/json"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/http/httputil"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager"
	types "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/keymanager/zk-snarks"
	usecases "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/key-manager/use-cases"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/service/formatters"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/services/key-manager/store"
)

const zksAccountPath = "/zk-snarks/accounts"

type ZKSController struct {
	vault    store.Vault
	useCases usecases.ZKSUseCases
}

func NewZKSController(vault store.Vault, useCases usecases.ZKSUseCases) *ZKSController {
	return &ZKSController{vault: vault, useCases: useCases}
}

// Add routes to router
func (c *ZKSController) Append(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/zk-snarks/namespaces").HandlerFunc(c.listNamespaces)
	router.Methods(http.MethodPost).Path(zksAccountPath).HandlerFunc(c.createAccount)
	router.Methods(http.MethodGet).Path(zksAccountPath).HandlerFunc(c.listAccounts)
	router.Methods(http.MethodGet).Path(zksAccountPath + "/{publicKey}").HandlerFunc(c.getAccount)
	router.Methods(http.MethodPost).Path(zksAccountPath + "/{publicKey}/sign").HandlerFunc(c.signPayload)
	router.Methods(http.MethodPost).Path(zksAccountPath + "/verify-signature").HandlerFunc(c.verifySignature)
}

// @Summary List zk-snarks Namespaces
// @Description List zk-snarks namespaces in the Vault
// @Produce json
// @Success 200 {object} []string "List of zk-snarks public namespaces"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /zk-snarks/namespaces [get]
func (c *ZKSController) listNamespaces(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	namespaces, err := c.vault.ZKSListNamespaces()
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(namespaces)
}

// @Summary Creates a new zk-snarks Account
// @Description Creates a new private key, stores it in the Vault and generates a public key given a chosen elliptic curve
// @Accept json
// @Produce json
// @Param request body zksnarks.CreateETHAccountRequest true "zk-snarks account creation request"
// @Success 200 {object} zksnarks.ETHAccountResponse "Created zk-snarks account"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /zk-snarks/accounts [post]
func (c *ZKSController) createAccount(rw http.ResponseWriter, request *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	accountRequest := &types.CreateZKSAccountRequest{}
	err := jsonutils.UnmarshalBody(request.Body, accountRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	accountResponse, err := c.vault.ZKSCreateAccount(accountRequest.Namespace)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(formatters.FormatZKSAccountResponse(accountResponse))
}

// @Summary List zk-snarks Accounts
// @Description List stored zk-snarks account in the Vault
// @Produce json
// @Param namespace query string false "namespace where key is stored"
// @Success 200 {object} []string "List of zk-snarks public accounts"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /zk-snarks/accounts [get]
func (c *ZKSController) listAccounts(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	namespace := req.URL.Query().Get("namespace")

	accountAddrs, err := c.vault.ZKSListAccounts(namespace)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_ = json.NewEncoder(rw).Encode(accountAddrs)
}

// @Summary Fetch zk-snarks Account
// @Description Get selected stored zk-snarks account in the Vault
// @Produce json
// @Param namespace query string false "namespace where key is stored"
// @Success 200 {object} zksnarks.ZKSAccountResponse "zk-snarks account"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /zk-snarks/accounts/{publicKey} [get]
func (c *ZKSController) getAccount(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	publicKey := mux.Vars(req)["publicKey"]
	namespace := req.URL.Query().Get("namespace")

	ethAcc, err := c.vault.ZKSGetAccount(publicKey, namespace)
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

// @Summary Signs an arbitrary message using an existing zk-snarks account
// @Description Signs an arbitrary message using EDDSA and the private key of an existing zk-snarks account
// @Accept json
// @Produce text/plain
// @Param request body keymanager.PayloadRequest true "Payload to sign"
// @Success 200 {string} string "Signed payload"
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 404 {object} httputil.ErrorResponse "Account not found"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /zk-snarks/accounts/{publicKey}/sign [post]
func (c *ZKSController) signPayload(rw http.ResponseWriter, req *http.Request) {
	signRequest := &keymanager.SignPayloadRequest{}
	err := jsonutils.UnmarshalBody(req.Body, signRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	publicKey := mux.Vars(req)["publicKey"]
	signature, err := c.vault.ZKSSign(publicKey, signRequest.Namespace, signRequest.Data)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	_, _ = rw.Write([]byte(signature))
}

// @Summary Verifies the signature of a message
// @Description Verifies if a message has been signed by the zk-snarks account passed as argument
// @Accept json
// @Param request body ethereum.VerifySignatureRequest true "signature and message to verify"
// @Success 204
// @Failure 400 {object} httputil.ErrorResponse "Invalid request"
// @Failure 422 {object} httputil.ErrorResponse "Failed to verify"
// @Failure 500 {object} httputil.ErrorResponse "Internal server error"
// @Router /zk-snarks/accounts/verify-signature [post]
func (c *ZKSController) verifySignature(rw http.ResponseWriter, request *http.Request) {
	verifyRequest := &types.VerifyPayloadRequest{}
	err := jsonutils.UnmarshalBody(request.Body, verifyRequest)
	if err != nil {
		httputil.WriteError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.useCases.VerifySignature().Execute(request.Context(), verifyRequest.PublicKey, verifyRequest.Signature, verifyRequest.Data)
	if err != nil {
		httputil.WriteHTTPErrorResponse(rw, err)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}
