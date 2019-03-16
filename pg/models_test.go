package pg

import (
	"context"
	"math/big"
	"testing"
	"time"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-pg/pg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/api/context-store.git/pg/migrations"
	common "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/common"
	ethereum "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/ethereum"
	trace "gitlab.com/ConsenSys/client/fr/core-stack/pkg.git/protos/trace"
)

type ContextStoreTestSuite struct {
	suite.Suite
	db *pg.DB
}

func (suite *ContextStoreTestSuite) SetupSuite() {
	// Create a test database
	db := pg.Connect(&pg.Options{
		Addr:     "127.0.0.1:5432",
		User:     "postgres",
		Password: "postgres",
		Database: "postgres",
	})

	testTable := "test"
	_, err := db.Exec(`DROP DATABASE IF EXISTS ?;`, pg.Q(testTable))
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE DATABASE ?;`, pg.Q(testTable))
	if err != nil {
		panic(err)
	}

	db.Close()

	// Create a connection to test database
	suite.db = pg.Connect(&pg.Options{
		Addr:     "127.0.0.1:5432",
		User:     "postgres",
		Password: "postgres",
		Database: "test",
	})
	migrations.Run(suite.db, "init")
}

func (suite *ContextStoreTestSuite) SetupTest() {
	oldVersion, newVersion, err := migrations.Run(suite.db, "up")
	if err != nil {
		suite.T().Errorf("Migrate up: %v", err)
	} else {
		suite.T().Logf("Migrated up from version=%v to version=%v", oldVersion, newVersion)
	}
}

func (suite *ContextStoreTestSuite) TearDownTest() {
	oldVersion, newVersion, err := migrations.Run(suite.db, "reset")
	if err != nil {
		suite.T().Errorf("Migrate down: %v", err)
	} else {
		suite.T().Logf("Migrated down from version=%v to version=%v", oldVersion, newVersion)
	}
}

func (suite *ContextStoreTestSuite) TearDownSuite() {
	// Close connection to test database
	suite.db.Close()

	// Drop test Database
	db := pg.Connect(&pg.Options{
		Addr:     "127.0.0.1:5432",
		User:     "postgres",
		Password: "postgres",
		Database: "postgres",
	})
	db.Exec(`DROP DATABASE test;`)
	db.Close()
}

type ModelsTestSuite struct {
	ContextStoreTestSuite
}

func (suite *ModelsTestSuite) TestStore() {
	s := traceStore{db: suite.db}

	txData := (&ethereum.TxData{}).
		SetNonce(10).
		SetTo(ethcommon.HexToAddress("0xAf84242d70aE9D268E2bE3616ED497BA28A7b62C")).
		SetValue(big.NewInt(100000)).
		SetGas(2000).
		SetGasPrice(big.NewInt(200000)).
		SetData(hexutil.MustDecode("0xabcd"))

	tr := &trace.Trace{
		Chain:    &common.Chain{Id: "0x3"},
		Metadata: &trace.Metadata{Id: "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11"},
		Tx: &ethereum.Transaction{
			TxData: txData,
			Raw:    "0xf86c0184ee6b280082529094ff778b716fc07d98839f48ddb88d8be583beb684872386f26fc1000082abcd29a0d1139ca4c70345d16e00f624622ac85458d450e238a48744f419f5345c5ce562a05bd43c512fcaf79e1756b2015fec966419d34d2a87d867b9618a48eca33a1a80",
			Hash:   "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210",
		},
	}

	status, storedAt, err := s.Store(context.Background(), tr)
	assert.Nil(suite.T(), err, "Should properly store trace")
	assert.Equal(suite.T(), "stored", status, "Default status should be correct")
	assert.True(suite.T(), time.Now().Sub(storedAt) < 5*time.Second, "Stored date should be close")

	_, _, err = s.Store(context.Background(), tr)
	assert.NotNil(suite.T(), err, "Unique constraint on TxHash should be violated")

	tr = &trace.Trace{}
	status, _, err = s.LoadByTxHash(context.Background(), "0x3", "0x0a0cafa26ca3f411e6629e9e02c53f23713b0033d7a72e534136104b5447a210", tr)
	assert.Nil(suite.T(), err, "Should properly store trace")
	assert.Equal(suite.T(), "stored", status, "Status should be correct")
	assert.Equal(suite.T(), "0x3", tr.GetChain().GetId(), "ChainID should be correct")
	assert.Equal(suite.T(), "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11", tr.GetMetadata().GetId(), "MetadataID should be correct")

	err = s.SetStatus(context.Background(), tr.GetMetadata().GetId(), "pending")
	assert.Nil(suite.T(), err, "Setting status to %q", "pending")

	status, sentAt, err := s.GetStatus(context.Background(), "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11")
	assert.Equal(suite.T(), "pending", status, "Status should be correct")
	assert.True(suite.T(), sentAt.Sub(storedAt) > 0, "Stored should be older than sent date")

	err = s.SetStatus(context.Background(), "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11", "error")
	assert.Nil(suite.T(), err, "Setting status to %q", "error")

	status, errorAt, err := s.GetStatus(context.Background(), "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11")
	assert.Equal(suite.T(), "error", status, "Status should be correct")
	assert.True(suite.T(), errorAt.Sub(sentAt) > 0, "Stored date should be close")

	status, _, err = s.LoadByTraceID(context.Background(), "a0ee-bc99-9c0b-4ef8-bb6d-6bb9-bd38-0a11", tr)
	assert.Equal(suite.T(), "error", status, "Status should be correct")
}

func TestMigrations(t *testing.T) {
	suite.Run(t, new(ModelsTestSuite))
}
