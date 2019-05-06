package hashicorp

// import (
// 	"testing"

// 	"github.com/stretchr/testify/suite"
// 	"gitlab.com/ConsenSys/client/fr/core-stack/infra/key-store.git/secretstore/testutils"
// )

// type HashicorpKeyStoreTestSuite struct {
// 	testutils.SecretStoreTestSuite
// }

// func (s *HashicorpKeyStoreTestSuite) SetupTest() {
// 	config := NewConfig()
// 	hashicorps, err := NewHashiCorp(config)
// 	if err != nil {
// 		panic(err)
// 	}
// 	s.Store = hashicorps
// }

// func TestMock(t *testing.T) {
// 	s := new(HashicorpKeyStoreTestSuite)
// 	suite.Run(t, s)
// }
