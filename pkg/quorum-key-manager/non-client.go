package quorumkeymanager

import (
	"context"

	"github.com/consensys/orchestrate/pkg/errors"
	qkm "github.com/consensys/quorum-key-manager/pkg/client"
	"github.com/consensys/quorum-key-manager/pkg/jsonrpc"
	aliastypes "github.com/consensys/quorum-key-manager/src/aliases/api/types"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
	utilstypes "github.com/consensys/quorum-key-manager/src/utils/api/types"
)

const errMessage = "Quorum Key Manager is disabled"

type NonClient struct{}

var _ qkm.KeyManagerClient = &NonClient{}

func NewNonClient() *NonClient {
	return &NonClient{}
}

func (n NonClient) SetSecret(ctx context.Context, storeName, id string, request *types.SetSecretRequest) (*types.SecretResponse, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) GetSecret(ctx context.Context, storeName, id, version string) (*types.SecretResponse, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) GetDeletedSecret(ctx context.Context, storeName, id string) (*types.SecretResponse, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) DeleteSecret(ctx context.Context, storeName, id string) error {
	return errors.DependencyFailureError(errMessage)
}

func (n NonClient) RestoreSecret(ctx context.Context, storeName, id string) error {
	return errors.DependencyFailureError(errMessage)
}

func (n NonClient) DestroySecret(ctx context.Context, storeName, id string) error {
	return errors.DependencyFailureError(errMessage)
}

func (n NonClient) ListSecrets(ctx context.Context, storeName string, limit, page uint64) ([]string, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) ListDeletedSecrets(ctx context.Context, storeName string, limit, page uint64) ([]string, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) CreateKey(ctx context.Context, storeName, id string, request *types.CreateKeyRequest) (*types.KeyResponse, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) ImportKey(ctx context.Context, storeName, id string, request *types.ImportKeyRequest) (*types.KeyResponse, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) SignKey(ctx context.Context, storeName, id string, request *types.SignBase64PayloadRequest) (string, error) {
	return "", errors.DependencyFailureError(errMessage)
}

func (n NonClient) GetKey(ctx context.Context, storeName, id string) (*types.KeyResponse, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) ListKeys(ctx context.Context, storeName string, limit, page uint64) ([]string, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) DeleteKey(ctx context.Context, storeName, id string) error {
	return errors.DependencyFailureError(errMessage)
}

func (n NonClient) GetDeletedKey(ctx context.Context, storeName, id string) (*types.KeyResponse, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) ListDeletedKeys(ctx context.Context, storeName string, limit, page uint64) ([]string, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) RestoreKey(ctx context.Context, storeName, id string) error {
	return errors.DependencyFailureError(errMessage)
}

func (n NonClient) DestroyKey(ctx context.Context, storeName, id string) error {
	return errors.DependencyFailureError(errMessage)
}

func (n NonClient) CreateEthAccount(ctx context.Context, storeName string, request *types.CreateEthAccountRequest) (*types.EthAccountResponse, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) ImportEthAccount(ctx context.Context, storeName string, request *types.ImportEthAccountRequest) (*types.EthAccountResponse, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) UpdateEthAccount(ctx context.Context, storeName, address string, request *types.UpdateEthAccountRequest) (*types.EthAccountResponse, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) SignMessage(ctx context.Context, storeName, account string, request *types.SignMessageRequest) (string, error) {
	return "", errors.DependencyFailureError(errMessage)
}

func (n NonClient) SignTypedData(ctx context.Context, storeName, address string, request *types.SignTypedDataRequest) (string, error) {
	return "", errors.DependencyFailureError(errMessage)
}

func (n NonClient) SignTransaction(ctx context.Context, storeName, address string, request *types.SignETHTransactionRequest) (string, error) {
	return "", errors.DependencyFailureError(errMessage)
}

func (n NonClient) SignQuorumPrivateTransaction(ctx context.Context, storeName, address string, request *types.SignQuorumPrivateTransactionRequest) (string, error) {
	return "", errors.DependencyFailureError(errMessage)
}

func (n NonClient) SignEEATransaction(ctx context.Context, storeName, address string, request *types.SignEEATransactionRequest) (string, error) {
	return "", errors.DependencyFailureError(errMessage)
}

func (n NonClient) GetEthAccount(ctx context.Context, storeName, address string) (*types.EthAccountResponse, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) ListEthAccounts(ctx context.Context, storeName string, limit, page uint64) ([]string, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) ListDeletedEthAccounts(ctx context.Context, storeName string, limit, page uint64) ([]string, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) DeleteEthAccount(ctx context.Context, storeName, address string) error {
	return errors.DependencyFailureError(errMessage)
}

func (n NonClient) DestroyEthAccount(ctx context.Context, storeName, address string) error {
	return errors.DependencyFailureError(errMessage)
}

func (n NonClient) RestoreEthAccount(ctx context.Context, storeName, address string) error {
	return errors.DependencyFailureError(errMessage)
}

func (n NonClient) VerifyKeySignature(ctx context.Context, request *utilstypes.VerifyKeySignatureRequest) error {
	return errors.DependencyFailureError(errMessage)
}

func (n NonClient) ECRecover(ctx context.Context, request *utilstypes.ECRecoverRequest) (string, error) {
	return "", errors.DependencyFailureError(errMessage)
}

func (n NonClient) VerifyMessage(ctx context.Context, request *utilstypes.VerifyRequest) error {
	return errors.DependencyFailureError(errMessage)
}

func (n NonClient) VerifyTypedData(ctx context.Context, request *utilstypes.VerifyTypedDataRequest) error {
	return errors.DependencyFailureError(errMessage)
}

func (n NonClient) CreateAlias(ctx context.Context, registry, aliasKey string, req *aliastypes.AliasRequest) (*aliastypes.AliasResponse, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) GetAlias(ctx context.Context, registry, aliasKey string) (*aliastypes.AliasResponse, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) UpdateAlias(ctx context.Context, registry, aliasKey string, req *aliastypes.AliasRequest) (*aliastypes.AliasResponse, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) DeleteAlias(ctx context.Context, registry, aliasKey string) error {
	return errors.DependencyFailureError(errMessage)
}

func (n NonClient) CreateRegistry(ctx context.Context, registry string, req *aliastypes.CreateRegistryRequest) (*aliastypes.RegistryResponse, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) GetRegistry(ctx context.Context, registry string) (*aliastypes.RegistryResponse, error) {
	return nil, errors.DependencyFailureError(errMessage)
}

func (n NonClient) DeleteRegistry(ctx context.Context, registry string) error {
	return errors.DependencyFailureError(errMessage)
}

func (n NonClient) Call(ctx context.Context, nodeID, method string, args ...interface{}) (*jsonrpc.ResponseMsg, error) {
	return nil, errors.DependencyFailureError(errMessage)
}
