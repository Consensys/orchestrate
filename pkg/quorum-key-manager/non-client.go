package quorumkeymanager

import (
	"context"

	"github.com/consensys/orchestrate/pkg/errors"
	qkm "github.com/consensys/quorum-key-manager/pkg/client"
	"github.com/consensys/quorum-key-manager/pkg/jsonrpc"
	types2 "github.com/consensys/quorum-key-manager/src/aliases/api/types"
	"github.com/consensys/quorum-key-manager/src/stores/api/types"
)

var _ qkm.KeyManagerClient = &NonClient{}

type NonClient struct {
}

func NewNonClient() *NonClient {
	return &NonClient{}
}

func (n NonClient) SetSecret(ctx context.Context, storeName, id string, request *types.SetSecretRequest) (*types.SecretResponse, error) {
	return nil, errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) GetSecret(ctx context.Context, storeName, id, version string) (*types.SecretResponse, error) {
	return nil, errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) GetDeletedSecret(ctx context.Context, storeName, id string) (*types.SecretResponse, error) {
	return nil, errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) DeleteSecret(ctx context.Context, storeName, id string) error {
	return errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) RestoreSecret(ctx context.Context, storeName, id string) error {
	return errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) DestroySecret(ctx context.Context, storeName, id string) error {
	return errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) ListSecrets(ctx context.Context, storeName string, limit, page uint64) ([]string, error) {
	return nil, errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) ListDeletedSecrets(ctx context.Context, storeName string, limit, page uint64) ([]string, error) {
	return nil, errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) CreateKey(ctx context.Context, storeName, id string, request *types.CreateKeyRequest) (*types.KeyResponse, error) {
	return nil, errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) ImportKey(ctx context.Context, storeName, id string, request *types.ImportKeyRequest) (*types.KeyResponse, error) {
	return nil, errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) SignKey(ctx context.Context, storeName, id string, request *types.SignBase64PayloadRequest) (string, error) {
	return "", errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) GetKey(ctx context.Context, storeName, id string) (*types.KeyResponse, error) {
	return nil, errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) ListKeys(ctx context.Context, storeName string, limit, page uint64) ([]string, error) {
	return nil, errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) DeleteKey(ctx context.Context, storeName, id string) error {
	return errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) GetDeletedKey(ctx context.Context, storeName, id string) (*types.KeyResponse, error) {
	return nil, errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) ListDeletedKeys(ctx context.Context, storeName string, limit, page uint64) ([]string, error) {
	return nil, errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) RestoreKey(ctx context.Context, storeName, id string) error {
	return errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) DestroyKey(ctx context.Context, storeName, id string) error {
	return errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) CreateEthAccount(ctx context.Context, storeName string, request *types.CreateEthAccountRequest) (*types.EthAccountResponse, error) {
	return nil, errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) ImportEthAccount(ctx context.Context, storeName string, request *types.ImportEthAccountRequest) (*types.EthAccountResponse, error) {
	return nil, errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) UpdateEthAccount(ctx context.Context, storeName, address string, request *types.UpdateEthAccountRequest) (*types.EthAccountResponse, error) {
	return nil, errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) SignMessage(ctx context.Context, storeName, account string, request *types.SignMessageRequest) (string, error) {
	return "", errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) SignTypedData(ctx context.Context, storeName, address string, request *types.SignTypedDataRequest) (string, error) {
	return "", errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) SignTransaction(ctx context.Context, storeName, address string, request *types.SignETHTransactionRequest) (string, error) {
	return "", errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) SignQuorumPrivateTransaction(ctx context.Context, storeName, address string, request *types.SignQuorumPrivateTransactionRequest) (string, error) {
	return "", errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) SignEEATransaction(ctx context.Context, storeName, address string, request *types.SignEEATransactionRequest) (string, error) {
	return "", errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) GetEthAccount(ctx context.Context, storeName, address string) (*types.EthAccountResponse, error) {
	return nil, errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) ListEthAccounts(ctx context.Context, storeName string, limit, page uint64) ([]string, error) {
	return nil, errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) ListDeletedEthAccounts(ctx context.Context, storeName string, limit, page uint64) ([]string, error) {
	return nil, errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) DeleteEthAccount(ctx context.Context, storeName, address string) error {
	return errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) DestroyEthAccount(ctx context.Context, storeName, address string) error {
	return errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) RestoreEthAccount(ctx context.Context, storeName, address string) error {
	return errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) VerifyKeySignature(ctx context.Context, request *types.VerifyKeySignatureRequest) error {
	return errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) ECRecover(ctx context.Context, request *types.ECRecoverRequest) (string, error) {
	return "", errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) VerifyMessage(ctx context.Context, request *types.VerifyRequest) error {
	return errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) VerifyTypedData(ctx context.Context, request *types.VerifyTypedDataRequest) error {
	return errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) CreateAlias(ctx context.Context, registry, aliasKey string, req types2.AliasRequest) (*types2.AliasResponse, error) {
	return nil, errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) GetAlias(ctx context.Context, registry, aliasKey string) (*types2.AliasResponse, error) {
	return nil, errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) UpdateAlias(ctx context.Context, registry, aliasKey string, req types2.AliasRequest) (*types2.AliasResponse, error) {
	return nil, errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) DeleteAlias(ctx context.Context, registry, aliasKey string) error {
	return errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) ListAliases(ctx context.Context, registry string) ([]types2.Alias, error) {
	return nil, errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) DeleteRegistry(ctx context.Context, registry string) error {
	return errors.DependencyFailureError("Quorum Key Manager is disabled")
}

func (n NonClient) Call(ctx context.Context, nodeID, method string, args ...interface{}) (*jsonrpc.ResponseMsg, error) {
	return nil, errors.DependencyFailureError("Quorum Key Manager is disabled")
}
