package base

import (
	"math/big"
	"sync"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/secretstore/mock"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
)

var testPKeys = []struct {
	prv string
	a   string
}{
	{"5FBB50BFF6DFAD35C4A374C9237BA2F7EAED9C6868E0108CB259B62D68029B1A", "0xdbb881a51CD4023E4400CEF3ef73046743f08da3"},
	{"86B021CCB810F26A30445B85F71E4C1596A11A97DDF9B9E348AC93D1DA6735BC", "0x6009608A02a7A15fd6689D6DaD560C44E9ab61Ff"},
	{"DD614C3B343E1B6DBD1B2811D4F146CC90337DEEF96AB97C353578E871B19D5E", "0xfF778b716FC07D98839f48DdB88D8bE583BEB684"},
	{"425D92F63A836F890F1690B34B6A25C2971EF8D035CD8EA8592FD1069BD151C6", "0xffbBa394DEf3Ff1df0941c6429887107f58d4e9b"},
}

var testChains = []struct {
	ID       string
	IsEIP155 bool
}{
	{"0x1ae3", true},
	{"0x3", true},
	{"0xbf6e", false},
}

func makeSignerInput(i int) (*common.Chain, ethcommon.Address, *ethtypes.Transaction) {
	chain := &common.Chain{
		Id:       testChains[i%len(testChains)].ID,
		IsEIP155: testChains[i%len(testChains)].IsEIP155,
	}
	address := ethcommon.HexToAddress(testPKeys[i%len(testPKeys)].a)
	tx := ethtypes.NewTransaction(
		10,
		ethcommon.HexToAddress("0xfF778b716FC07D98839f48DdB88D8bE583BEB684"),
		big.NewInt(1000),
		100,
		big.NewInt(1000),
		hexutil.MustDecode("0xa2bcdef3"),
	)
	return chain, address, tx
}

// BaseKeyStoreTestSuite is a test suit for TraceStore
type BaseKeyStoreTestSuite struct {
	suite.Suite
	Store *KeyStore
}

func (suite *BaseKeyStoreTestSuite) SetupTest() {
	suite.Store = NewKeyStore(mock.NewSecretStore())
	for _, priv := range testPKeys {
		suite.Store.ImportPrivateKey(priv.prv)
	}
}

func (suite *BaseKeyStoreTestSuite) TestKeyStore() {
	// Feed input channel and then close it
	rounds := 1000
	wg := &sync.WaitGroup{}
	out := make(chan []byte, rounds)
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			raw, _, _ := suite.Store.SignTx(makeSignerInput(i))
			out <- raw
		}(i)
	}
	wg.Wait()
	close(out)

	assert.Len(suite.T(), out, rounds, "Count of signatures should be correct")
	for raw := range out {
		assert.True(suite.T(), len(raw) > 95, "Expected transaction to be signed but got %q", hexutil.Encode(raw))
	}
}

func TestKeyStore(t *testing.T) {
	s := new(BaseKeyStoreTestSuite)
	suite.Run(t, s)
}

//TestSecretStore must be run along with a vault container in development mode
//It will sequentially writes a secret, list all the secrets, get the secret then delete it.
// func TestKeyStore(t *testing.T) {
// 	config := hashicorp.NewConfig()
// 	hashicorpsSS, err := hashicorp.NewHashicorps(config)
// 	if err != nil {
// 		t.Errorf("Error when instantiating the vault : %v", err.Error())
// 	}

// 	err = hashicorpsSS.InitVault()
// 	if err != nil {
// 		t.Errorf("Error initializing the vault : %v", err.Error())
// 	}

// 	keystore := NewKeyStore(hashicorpsSS)

// 	_, err = keystore.GenerateWallet()
// 	if err != nil {
// 		t.Errorf("Error while generating a new wallet : %v", err.Error())
// 	}
// }
