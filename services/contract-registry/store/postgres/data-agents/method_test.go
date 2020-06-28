// +build unit
// +build !race
// +build !integration

package dataagents

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres"
	pgTestUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres/testutils"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/common"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/models"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/postgres/migrations"
)

type methodTestSuite struct {
	suite.Suite
	dataagent  *PGMethod
	codeHashDA *PGCodeHash
	pg         *pgTestUtils.PGTestHelper
}

var (
	methodSig0 = "methodSig0()"
	methodSig1 = "methodSig1()"
)

func TestPGMethod(t *testing.T) {
	s := new(methodTestSuite)
	suite.Run(t, s)
}

func (s *methodTestSuite) SetupSuite() {
	s.pg, _ = pgTestUtils.NewPGTestHelper(nil, migrations.Collection)
	s.pg.InitTestDB(s.T())
}

func (s *methodTestSuite) SetupTest() {
	s.pg.UpgradeTestDB(s.T())
	s.codeHashDA = NewPGCodeHash(s.pg.DB)
	s.dataagent = NewPGMethod(s.pg.DB)
}

func (s *methodTestSuite) TearDownTest() {
	s.pg.DowngradeTestDB(s.T())
}

func (s *methodTestSuite) TearDownSuite() {
	s.pg.DropTestDB(s.T())
}

func (s *methodTestSuite) TestPGMethod_InsertMultipleTx() {
	selector4 := new([4]byte)
	selector := crypto.Keccak256([]byte(methodSig0))[:4]
	copy(selector4[:], selector)

	s.T().Run("should insert multiple models successfully", func(t *testing.T) {
		methods := []*models.MethodModel{
			{
				Codehash: "codeHash",
				Selector: *selector4,
				ABI:      "ABI",
			},
			{
				Codehash: "codeHash",
				Selector: *selector4,
				ABI:      "ABI",
			},
		}
		err := s.dataagent.InsertMultiple(context.Background(), &methods)

		assert.NoError(t, err)
		assert.Equal(t, 1, methods[0].ID)
		assert.Equal(t, 2, methods[1].ID)
	})

	s.T().Run("should insert multiple models successfully in TX", func(t *testing.T) {
		tx, _ := s.pg.DB.Begin()
		ctx := postgres.WithTx(context.Background(), tx)

		methods := []*models.MethodModel{
			{
				Codehash: "codeHash",
				Selector: *selector4,
				ABI:      "ABI",
			},
		}
		err := s.dataagent.InsertMultiple(ctx, &methods)
		_ = tx.Commit()

		assert.NoError(t, err)
		assert.NotEmpty(t, methods[0].ID)
	})

	s.T().Run("should return PostgresConnectionError if insert fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)

		methods := []*models.MethodModel{
			{
				Codehash: "codeHash",
				Selector: *selector4,
				ABI:      "ABI",
			},
		}
		err := s.dataagent.InsertMultiple(context.Background(), &methods)

		assert.True(t, errors.IsPostgresConnectionError(err))

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}

func (s *methodTestSuite) TestPGMethod_FindOneByAccountAndSelector() {
	s.T().Run("should find one method successfully", func(t *testing.T) {
		s.insertMethods()

		selector := crypto.Keccak256([]byte(methodSig0))[:4]
		method, err := s.dataagent.FindOneByAccountAndSelector(context.Background(), &common.AccountInstance{
			ChainId: "chainId",
			Account: "address",
		}, selector)
		assert.NoError(t, err)
		assert.Equal(t, method.ABI, "ABI")

		selector1 := crypto.Keccak256([]byte(methodSig1))[:4]
		method, err = s.dataagent.FindOneByAccountAndSelector(context.Background(), &common.AccountInstance{
			ChainId: "chainId",
			Account: "address",
		}, selector1)
		assert.NoError(t, err)
		assert.Equal(t, method.ABI, "ABI1")
	})

	s.T().Run("should return NotFoundError if no event is found", func(t *testing.T) {
		selector := crypto.Keccak256([]byte("unknown"))[:4]
		_, err := s.dataagent.FindOneByAccountAndSelector(context.Background(), &common.AccountInstance{
			ChainId: "unknown",
			Account: "unknown",
		}, selector)
		assert.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should return PostgresConnectionError if find fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)

		selector := crypto.Keccak256([]byte(methodSig0))[:4]
		_, err := s.dataagent.FindOneByAccountAndSelector(context.Background(), &common.AccountInstance{
			ChainId: "chainId",
			Account: "address",
		}, selector)
		assert.True(t, errors.IsPostgresConnectionError(err))

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}

func (s *methodTestSuite) TestPGMethod_FindDefaultBySelector() {
	s.T().Run("should find default events successfully", func(t *testing.T) {
		s.insertMethods()

		selector := crypto.Keccak256([]byte(methodSig0))[:4]
		defaultEvents, err := s.dataagent.FindDefaultBySelector(context.Background(), selector)

		assert.NoError(t, err)
		assert.Len(t, defaultEvents, 1)
	})

	s.T().Run("should return NotFoundError if no default events are found", func(t *testing.T) {
		selector := crypto.Keccak256([]byte("unknown"))[:4]
		_, err := s.dataagent.FindDefaultBySelector(context.Background(), selector)
		assert.True(t, errors.IsNotFoundError(err))
	})

	s.T().Run("should return PostgresConnectionError if find fails", func(t *testing.T) {
		// We drop the DB to make the test fail
		s.pg.DropTestDB(t)

		selector := crypto.Keccak256([]byte(methodSig0))[:4]
		_, err := s.dataagent.FindDefaultBySelector(context.Background(), selector)
		assert.True(t, errors.IsPostgresConnectionError(err))

		// We bring it back up
		s.pg.InitTestDB(t)
	})
}

func (s *methodTestSuite) insertMethods() {
	_ = s.codeHashDA.Insert(context.Background(), &models.CodehashModel{
		ChainID:  "chainId",
		Address:  "address",
		Codehash: "codeHash",
	})

	_ = s.codeHashDA.Insert(context.Background(), &models.CodehashModel{
		ChainID:  "chainId",
		Address:  "address1",
		Codehash: "codeHash1",
	})

	selector04 := new([4]byte)
	selector0 := crypto.Keccak256([]byte(methodSig0))[:4]
	copy(selector04[:], selector0)

	selector14 := new([4]byte)
	selector1 := crypto.Keccak256([]byte(methodSig1))[:4]
	copy(selector14[:], selector1)

	methods := []*models.MethodModel{
		{
			Codehash: "codeHash",
			Selector: *selector04,
			ABI:      "ABI",
		},
		{
			Codehash: "codeHash",
			Selector: *selector14,
			ABI:      "ABI1",
		},
	}

	_ = s.dataagent.InsertMultiple(context.Background(), &methods)
}
