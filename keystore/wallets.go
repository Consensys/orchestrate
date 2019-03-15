package keystore

import(
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/aws-secret-manager.git/secretstore"
	"github.com/ethereum/go-ethereum/common"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/crypto"
	"encoding/hex"
	"fmt"
)

// Wallet is a container
type Wallet struct {
	address common.Address
	priv *ecdsa.PrivateKey
	sec secretstore.SecretStore
}

// NewWallet construct a wallet object
func NewWallet(sec secretstore.SecretStore) *Wallet {
	return &Wallet{
		sec: sec,
	}
}

// GenerateWallet create a keypair
func (wal *Wallet) Generate() (err error) {

	wal = &Wallet{}
	
	wal.priv, err = crypto.GenerateKey()
	if err != nil {
		return err
	}

	pub := wal.priv.PublicKey
	wal.address = crypto.PubkeyToAddress(pub)
	return nil
}

// Store saves the wallet in the vault
func (wal *Wallet) Store() (err error) {
	return wal.sec.Store(
		wal.address.Hex(),
		hex.EncodeToString(crypto.FromECDSA(wal.priv)),
	)

}

// GetWallet returns a wallet object from an address if its stored in the vault
func (wal *Wallet) Load(a *common.Address) (err error) {

	wal = &Wallet{
		address: *a,
	}

	priv, err := wal.sec.Load(a.Hex())
	if err != nil {
		return err
	}

	wal.priv, err = crypto.HexToECDSA(priv)
	if err != nil {
		return fmt.Errorf("Could not deserialize %v...%v", priv[:5], priv[len(priv)-5:])
	}

	return nil

}

// GetPriv returns the private key of Wallet
func (wal *Wallet) GetPriv() (*ecdsa.PrivateKey) {
	return wal.priv
}

//GetAddress returns the address of the wallet
func (wal *Wallet) GetAddress() (*common.Address) {
	return &wal.address
}