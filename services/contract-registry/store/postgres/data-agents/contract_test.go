// +build unit
// +build !race
// +build !integration

package dataagents

// import (
// 	"context"
// 	"fmt"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// 	pgTestUtils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/database/postgres/testutils"
// 	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/errors"
// 	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store"
// 	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/mocks"
// 	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/models"
// 	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/contract-registry/store/postgres/migrations"
// )

// type contractTestSuite struct {
// 	suite.Suite
// 	dataagent        store.ContractDataAgent
// 	mockRepositoryDA *mocks.MockRepositoryDataAgent
// 	mockArtifactDA   *mocks.MockArtifactDataAgent
// 	mockTagDA        *mocks.MockTagDataAgent
// 	mockMethodDA     *mocks.MockMethodDataAgent
// 	mockEventDA      *mocks.MockEventDataAgent
// 	pg               *pgTestUtils.PGTestHelper
// }

// func TestPGContract(t *testing.T) {
// 	s := new(contractTestSuite)
// 	suite.Run(t, s)
// }

// func (s *contractTestSuite) SetupSuite() {
// 	s.pg , _ = pgTestUtils.NewPGTestHelper(nil, migrations.Collection)
// 	s.pg.InitTestDB(s.T())
// }

// func (s *contractTestSuite) SetupTest() {
// 	ctrl := gomock.NewController(s.T())
// 	defer ctrl.Finish()

// 	s.pg.UpgradeTestDB(s.T())

// 	s.mockRepositoryDA = mocks.NewMockRepositoryDataAgent(ctrl)
// 	s.mockArtifactDA = mocks.NewMockArtifactDataAgent(ctrl)
// 	s.mockTagDA = mocks.NewMockTagDataAgent(ctrl)
// 	s.mockMethodDA = mocks.NewMockMethodDataAgent(ctrl)
// 	s.mockEventDA = mocks.NewMockEventDataAgent(ctrl)

// 	s.dataagent = NewPGContract(s.pg.DB, s.mockRepositoryDA, s.mockArtifactDA, s.mockTagDA, s.mockMethodDA, s.mockEventDA)
// }

// func (s *contractTestSuite) TearDownTest() {
// 	s.pg.DowngradeTestDB(s.T())
// }

// func (s *contractTestSuite) TearDownSuite() {
// 	s.pg.DropTestDB(s.T())
// }

// func (s *contractTestSuite) TestPGContract_Insert() {
// 	methods := []*models.MethodModel{
// 		{
// 			Codehash: "codeHash",
// 			Selector: [4]byte{58, 58},
// 			ABI:      "ABI",
// 		},
// 	}
// 	events := []*models.EventModel{
// 		{
// 			Codehash:          "codeHash",
// 			SigHash:           "sigHash",
// 			IndexedInputCount: 0,
// 			ABI:               "ABI",
// 		},
// 	}
// 	dataAgentError := fmt.Errorf("error")

// 	s.T().Run("should insert contract successfully", func(t *testing.T) {
// 		s.mockRepositoryDA.EXPECT().SelectOrInsert(gomock.Any(), gomock.Any()).Return(nil)
// 		s.mockArtifactDA.EXPECT().SelectOrInsert(gomock.Any(), gomock.Any()).Return(nil)
// 		s.mockTagDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
// 		s.mockMethodDA.EXPECT().InsertMultiple(gomock.Any(), gomock.Any()).Return(nil)
// 		s.mockEventDA.EXPECT().InsertMultiple(gomock.Any(), gomock.Any()).Return(nil)

// 		err := s.dataagent.Insert(
// 			context.Background(),
// 			"name",
// 			"tag",
// 			"abi",
// 			"bytecode",
// 			"deployedBytecode",
// 			"codeHash",
// 			&methods,
// 			&events,
// 		)

// 		assert.Nil(t, err)
// 	})

// 	s.T().Run("should insert contract with empty methods successfully", func(t *testing.T) {
// 		s.mockRepositoryDA.EXPECT().SelectOrInsert(gomock.Any(), gomock.Any()).Return(nil)
// 		s.mockArtifactDA.EXPECT().SelectOrInsert(gomock.Any(), gomock.Any()).Return(nil)
// 		s.mockTagDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)

// 		emptyMethods := []*models.MethodModel{}
// 		emptyEvents := []*models.EventModel{}

// 		err := s.dataagent.Insert(
// 			context.Background(),
// 			"name",
// 			"tag",
// 			"abi",
// 			"bytecode",
// 			"deployedBytecode",
// 			"codeHash",
// 			&emptyMethods,
// 			&emptyEvents,
// 		)

// 		assert.Nil(t, err)
// 	})

// 	s.T().Run("should fail if repository data agent fails", func(t *testing.T) {
// 		s.mockRepositoryDA.EXPECT().SelectOrInsert(gomock.Any(), gomock.Any()).Return(dataAgentError)

// 		err := s.dataagent.Insert(
// 			context.Background(),
// 			"name",
// 			"tag",
// 			"abi",
// 			"bytecode",
// 			"deployedBytecode",
// 			"codeHash",
// 			&methods,
// 			&events,
// 		)

// 		assert.Equal(t, errors.FromError(dataAgentError).ExtendComponent(component), err)
// 	})

// 	s.T().Run("should fail if artifact data agent fails", func(t *testing.T) {
// 		s.mockRepositoryDA.EXPECT().SelectOrInsert(gomock.Any(), gomock.Any()).Return(nil)
// 		s.mockArtifactDA.EXPECT().SelectOrInsert(gomock.Any(), gomock.Any()).Return(dataAgentError)

// 		err := s.dataagent.Insert(
// 			context.Background(),
// 			"name",
// 			"tag",
// 			"abi",
// 			"bytecode",
// 			"deployedBytecode",
// 			"codeHash",
// 			&methods,
// 			&events,
// 		)
// 		assert.Equal(t, errors.FromError(dataAgentError).ExtendComponent(component), err)
// 	})

// 	s.T().Run("should fail if tag data agent fails", func(t *testing.T) {
// 		s.mockRepositoryDA.EXPECT().SelectOrInsert(gomock.Any(), gomock.Any()).Return(nil)
// 		s.mockArtifactDA.EXPECT().SelectOrInsert(gomock.Any(), gomock.Any()).Return(nil)
// 		s.mockTagDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(dataAgentError)

// 		err := s.dataagent.Insert(
// 			context.Background(),
// 			"name",
// 			"tag",
// 			"abi",
// 			"bytecode",
// 			"deployedBytecode",
// 			"codeHash",
// 			&methods,
// 			&events,
// 		)
// 		assert.Equal(t, errors.FromError(dataAgentError).ExtendComponent(component), err)
// 	})

// 	s.T().Run("should fail if method data agent fails", func(t *testing.T) {
// 		s.mockRepositoryDA.EXPECT().SelectOrInsert(gomock.Any(), gomock.Any()).Return(nil)
// 		s.mockArtifactDA.EXPECT().SelectOrInsert(gomock.Any(), gomock.Any()).Return(nil)
// 		s.mockTagDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
// 		s.mockMethodDA.EXPECT().InsertMultiple(gomock.Any(), gomock.Any()).Return(dataAgentError)

// 		err := s.dataagent.Insert(
// 			context.Background(),
// 			"name",
// 			"tag",
// 			"abi",
// 			"bytecode",
// 			"deployedBytecode",
// 			"codeHash",
// 			&methods,
// 			&events,
// 		)
// 		assert.Equal(t, errors.FromError(dataAgentError).ExtendComponent(component), err)
// 	})

// 	s.T().Run("should fail if event data agent fails", func(t *testing.T) {
// 		s.mockRepositoryDA.EXPECT().SelectOrInsert(gomock.Any(), gomock.Any()).Return(nil)
// 		s.mockArtifactDA.EXPECT().SelectOrInsert(gomock.Any(), gomock.Any()).Return(nil)
// 		s.mockTagDA.EXPECT().Insert(gomock.Any(), gomock.Any()).Return(nil)
// 		s.mockMethodDA.EXPECT().InsertMultiple(gomock.Any(), gomock.Any()).Return(nil)
// 		s.mockEventDA.EXPECT().InsertMultiple(gomock.Any(), gomock.Any()).Return(dataAgentError)

// 		err := s.dataagent.Insert(
// 			context.Background(),
// 			"name",
// 			"tag",
// 			"abi",
// 			"bytecode",
// 			"deployedBytecode",
// 			"codeHash",
// 			&methods,
// 			&events,
// 		)
// 		assert.Equal(t, errors.FromError(dataAgentError).ExtendComponent(component), err)
// 	})

// 	s.T().Run("should fail if DB connection fails", func(t *testing.T) {
// 		s.pg.DropTestDB(t)

// 		err := s.dataagent.Insert(
// 			context.Background(),
// 			"name",
// 			"tag",
// 			"abi",
// 			"bytecode",
// 			"deployedBytecode",
// 			"codeHash",
// 			&methods,
// 			&events,
// 		)

// 		assert.True(t, errors.IsPostgresConnectionError(err))

// 		s.pg.InitTestDB(t)
// 	})
// }
