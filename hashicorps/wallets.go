package hashicorps

import(
	"github.com/hashicorp/vault/api"
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
	pub *ecdsa.PublicKey
	secretStr string
	sec *Secret
}

// GenerateWallet create a keypair
func GenerateWallet() (wal *Wallet, err error) {

	wal = &Wallet{}
	
	wal.priv, err = crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	pub := wal.priv.PublicKey
	wal.address = crypto.PubkeyToAddress(pub)
	return wal, nil
}

// Store saves the wallet in the vault
func (wal *Wallet) Store(client *api.Client) (err error) {

	wal.sec = NewSecret().
		SetKey(wal.address.Hex()).
		SetValue(hex.EncodeToString(crypto.FromECDSA(wal.priv))).
		SetClient(client)

	_, err = wal.sec.SaveNew()
	if err != nil {
		return err
	}
	return nil
}

// GetWallet returns a wallet object from an address if its stored in the vault
func GetWallet(client *api.Client, a *common.Address) (wal *Wallet, err error) {

	wal = &Wallet{}
	wal.sec = SecretFromKey(a.Hex())
	wal.sec.SetClient(client)

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

// GetPriv returns the private key of Wallet
func (wal *Wallet) GetPriv() (*ecdsa.PrivateKey) {
	return wal.priv
}

//GetAddress returns the address of the wallet
func (wal *Wallet) GetAddress() (*common.Address) {
	return &wal.address
}