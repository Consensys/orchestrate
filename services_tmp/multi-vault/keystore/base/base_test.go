package base

import (
	"math/big"
	"sync"
	"testing"

	"github.com/ConsenSys/golang-utils/ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/types/chain"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/ethereum.git/types"
	"gitlab.com/ConsenSys/client/fr/core-stack/service/multi-vault.git/secretstore/mock"
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

var testChainsIds = []int64{
	10,
	3,
	13,
}

var arbitraryMsg = []string{
	"This is not a very long message to hash",
	"This is a bit longer but it does'nt tells a lot. So I think we should write some more text",
	"Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old.",
	"The ate pairing and its variations are simply optimized versions of the Tate pairing when restricted to the eigenspaces of Frobenius. Denote with π q the Frobenius endomorphism, i.e. π q : E → E : (x, y) 7→ (x q , y q ) and define G 1 = E[r] ∩ Ker(π q − [1]) = E(F q )[r] and G 2 = E[r] ∩ Ker(π q − [q]).",
	"Rust is for people who crave speed and stability in a language. By speed, we mean the speed of the programs that you can create with Rust and the speed at which Rust lets you write them. The Rust compiler’s checks ensure stability through feature additions and refactoring. This is in contrast to the brittle legacy code in languages without these checks, which developers are often afraid to modify. By striving for zero-cost abstractions, higher-level features that compile to lower-level code as fast as code written manually, Rust endeavors to make safe code be fast code as well.",
}

func makeSignTxInput(i int) (*chain.Chain, ethcommon.Address, *ethtypes.Transaction) {
	netChain := &chain.Chain{
		Id: big.NewInt(testChainsIds[i%len(testChainsIds)]).Bytes(),
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
	return netChain, address, tx
}

func makeSignMsgInput(i int) (a ethcommon.Address, msg string) {
	return ethcommon.HexToAddress(testPKeys[i%len(testPKeys)].a),
		arbitraryMsg[i%len(arbitraryMsg)]
}

// BaseKeyStoreTestSuite is a test suit for TraceStore
type BaseKeyStoreTestSuite struct {
	suite.Suite
	Store *KeyStore
}

func (s *BaseKeyStoreTestSuite) SetupTest() {
	s.Store = NewKeyStore(mock.NewSecretStore())
}

// TestSignTx is a test suit for KeyStore that test ethereum signature
func (s *BaseKeyStoreTestSuite) TestSignTx() {
	for _, priv := range testPKeys {
		err := s.Store.ImportPrivateKey(priv.prv)
		assert.NoError(s.T(), err)
	}

	// Feed input channel and then close it
	rounds := 1000
	wg := &sync.WaitGroup{}
	out := make(chan []byte, rounds)
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			raw, _, _ := s.Store.SignTx(makeSignTxInput(i))
			out <- raw
		}(i)
	}
	wg.Wait()
	close(out)

	assert.Len(s.T(), out, rounds, "Count of signatures should be correct")
	for raw := range out {
		assert.True(s.T(), len(raw) > 95, "Expected transaction to be signed but got %q", hexutil.Encode(raw))
	}
}

func (s *BaseKeyStoreTestSuite) TestSignMsg() {
	for _, priv := range testPKeys {
		err := s.Store.ImportPrivateKey(priv.prv)
		assert.Nil(s.T(), err)
	}

	// Feed input channel and then close it
	rounds := 20
	for i := 0; i < rounds; i++ {
		address, msg := makeSignMsgInput(i)
		signature, hash, err := s.Store.SignMsg(address, msg)
		assert.Nil(s.T(), err, "The msg has not been signed")

		recoveredAddress, err := ethereum.EcRecover(*hash, signature)
		assert.Nil(s.T(), err)
		assert.Equal(s.T(), address, recoveredAddress)
	}
}

func (s *BaseKeyStoreTestSuite) TestGenerateWallet() {
	_, err := s.Store.GenerateWallet()
	assert.Nil(s.T(), err, "Wallet should be generated")
}

func TestKeyStore(t *testing.T) {
	s := new(BaseKeyStoreTestSuite)
	suite.Run(t, s)
}

var privateKey = "8f2a55949038a9610f50fb23b5883af3b4ecb3c3bb792cbcefbd1542c692be63"

func TestPrivateTxSigning(t *testing.T) {
	expected := "0xf90268808203e8832dc6c094000000000000000000000000000000000000000080b901cb608060405234801561001057600080fd5b5060008054600160a060020a03191633179055610199806100326000396000f3fe6080604052600436106100565763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416633fa4f245811461005b5780636057361d1461008257806367e404ce146100ae575b600080fd5b34801561006757600080fd5b506100706100ec565b60408051918252519081900360200190f35b34801561008e57600080fd5b506100ac600480360360208110156100a557600080fd5b50356100f2565b005b3480156100ba57600080fd5b506100c3610151565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b60025490565b604080513381526020810183905281517fc9db20adedc6cf2b5d25252b101ab03e124902a73fcb12b753f3d1aaa2d8f9f5929181900390910190a16002556001805473ffffffffffffffffffffffffffffffffffffffff191633179055565b60015473ffffffffffffffffffffffffffffffffffffffff169056fea165627a7a72305820c7f729cb24e05c221f5aa913700793994656f233fe2ce3b9fd9a505ea17e8d8a00297ca0e4d0616c956ce5119719c1c4890ee8d4ac77595fdf6780878b0c3471834f29d6a0409fdaff39f7e85ac2a22ec0e5966a1b4673d187cff3cd6a34dbf5735cef7c42ac41316156744d784c4355486d425648586f5a7a7a42675062572f776a3561784470573958386c393153476f3dc08a72657374726963746564"
	data := "0x608060405234801561001057600080fd5b5060008054600160a060020a03191633179055610199806100326000396000f3fe6080604052600436106100565763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416633fa4f245811461005b5780636057361d1461008257806367e404ce146100ae575b600080fd5b34801561006757600080fd5b506100706100ec565b60408051918252519081900360200190f35b34801561008e57600080fd5b506100ac600480360360208110156100a557600080fd5b50356100f2565b005b3480156100ba57600080fd5b506100c3610151565b6040805173ffffffffffffffffffffffffffffffffffffffff9092168252519081900360200190f35b60025490565b604080513381526020810183905281517fc9db20adedc6cf2b5d25252b101ab03e124902a73fcb12b753f3d1aaa2d8f9f5929181900390910190a16002556001805473ffffffffffffffffffffffffffffffffffffffff191633179055565b60015473ffffffffffffffffffffffffffffffffffffffff169056fea165627a7a72305820c7f729cb24e05c221f5aa913700793994656f233fe2ce3b9fd9a505ea17e8d8a0029"

	privateFrom := "A1aVtMxLCUHmBVHXoZzzBgPbW/wj5axDpW9X8l91SGo="

	store := NewKeyStore(mock.NewSecretStore())
	err := store.ImportPrivateKey(privateKey)
	assert.NoError(t, err)

	chainID := 44
	tx := ethtypes.NewTransaction(
		0,
		ethcommon.HexToAddress("0x0"),
		big.NewInt(0),
		3000000,
		big.NewInt(1000),
		hexutil.MustDecode(data),
	)

	privateArgs := &types.PrivateArgs{
		PrivateFrom:   privateFrom,
		PrivateFor:    []string{},
		PrivateTxType: "restricted",
	}

	bytes, _, err := store.SignPrivateEEATx(chain.FromInt(int64(chainID)), ethcommon.HexToAddress("0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73"), tx, privateArgs)
	assert.NoError(t, err)
	assert.Equal(t, expected, hexutil.Encode(bytes))
}

func TestPrivateTesseraTxSigning(t *testing.T) {
	expected := "0xf865808203e8832dc6c094000000000000000000000000000000000000000080831234561ba0e04e1c11a8626c77fc2b61c246b066f765613b6104a824bb229889c45dae9922a018b07c72e29e7869422680979e5a4dfe0fd9fbcaa62bb3d0c9320a30b03d91c4"
	data := "0x123456"

	store := NewKeyStore(mock.NewSecretStore())
	err := store.ImportPrivateKey(privateKey)
	assert.NoError(t, err)

	chainID := 44
	tx := ethtypes.NewTransaction(
		0,
		ethcommon.HexToAddress("0x0"),
		big.NewInt(0),
		3000000,
		big.NewInt(1000),
		hexutil.MustDecode(data),
	)

	bytes, _, err := store.SignPrivateTesseraTx(chain.FromInt(int64(chainID)), ethcommon.HexToAddress("0xFE3B557E8Fb62b89F4916B721be55cEb828dBd73"), tx)
	assert.NoError(t, err)
	assert.Equal(t, expected, hexutil.Encode(bytes))
}
