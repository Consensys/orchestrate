package aws

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
	"encoding/hex"
	"fmt"
)

// Wallet is a container for private keys
type Wallet struct {
	address common.Address
	priv *ecdsa.PrivateKey
	pub *ecdsa.PublicKey
	secretStr string
	sec *Secret
}

// EmptyWallet is the default constructor of Wallet
func EmptyWallet() (*Wallet) {
	return &Wallet{}
}

// GenerateWallet returns a generated and saved wallet
func GenerateWallet(client *secretsmanager.SecretsManager) (wal *Wallet, err error) {

	wal = EmptyWallet() 
	wal.priv, err = crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	pub := wal.priv.PublicKey
	wal.address = crypto.PubkeyToAddress(pub)

	wal.sec = NewSecret().
		SetKey(wal.address.Hex()).
		SetValue(hex.EncodeToString(crypto.FromECDSA(wal.priv))).
		SetClient(client)

	_, err = wal.sec.SaveNew()
	if err != nil {
		return nil, err
	}
	return wal, err
}

// GetWallet construct a Wallet with the corresponding address
func GetWallet(client *secretsmanager.SecretsManager, a *common.Address) (wal *Wallet, err error) {
	
	wal = EmptyWallet()
	wal.sec, err = SecretFromKey(a.Hex())

	_, err = wal.sec.GetValue()
	if err != nil {
		return nil, err
	}

	wal.priv, err = crypto.HexToECDSA(wal.sec.value)
	if err != nil {
		return nil, fmt.Errorf("Could not deserialize %v", wal.sec.value)
	}

	return wal, nil
}




