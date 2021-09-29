// +build unit

package chainregistry

import (
	"context"
	"github.com/consensys/orchestrate/pkg/sdk/client/mock"
	"github.com/consensys/orchestrate/pkg/types/api"
	"github.com/consensys/orchestrate/pkg/types/testutils"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/consensys/orchestrate/services/tx-listener/dynamic"
	"github.com/consensys/orchestrate/services/tx-listener/session/ethereum/offset"
)

var mockChain = testutils.FakeChainResponse()

type ManagerTestSuite struct {
	suite.Suite
	Manager offset.Manager
	client  *mock.MockOrchestrateClient
}

func (s *ManagerTestSuite) SetupTest() {
	ctrl := gomock.NewController(s.T())
	s.client = mock.NewMockOrchestrateClient(ctrl)

	s.Manager = NewManager(s.client)
}

func (s *ManagerTestSuite) TestManagerLastBlock() {
	updatedCurrentBlock := uint64(1)
	chain := &dynamic.Chain{
		UUID: mockChain.UUID,
		Listener: dynamic.Listener{
			CurrentBlock: 0,
		},
	}

	s.client.EXPECT().GetChain(gomock.Any(), chain.UUID).Return(mockChain, nil)
	s.client.EXPECT().UpdateChain(gomock.Any(), chain.UUID, &api.UpdateChainRequest{Listener: &api.UpdateListenerRequest{
		CurrentBlock: updatedCurrentBlock,
	}}).Return(mockChain, nil)

	lastBlockNumber, err := s.Manager.GetLastBlockNumber(context.Background(), chain)
	assert.NoError(s.T(), err, "GetLastBlockNumber should not error")
	assert.Equal(s.T(), mockChain.ListenerCurrentBlock, lastBlockNumber, "Lastblock should be correct")

	err = s.Manager.SetLastBlockNumber(context.Background(), chain, updatedCurrentBlock)
	assert.NoError(s.T(), err, "SetLastBlockNumber should not error")
}

func (s *ManagerTestSuite) TestManagerLastBlock_ignored() {
	updatedCurrentBlock := uint64(0)
	chain := &dynamic.Chain{
		UUID: mockChain.UUID,
		Listener: dynamic.Listener{
			CurrentBlock: 0,
		},
	}

	s.client.EXPECT().GetChain(gomock.Any(), chain.UUID).Return(mockChain, nil)

	lastBlockNumber, err := s.Manager.GetLastBlockNumber(context.Background(), chain)
	assert.NoError(s.T(), err, "GetLastBlockNumber should not error")
	assert.Equal(s.T(), mockChain.ListenerCurrentBlock, lastBlockNumber, "Lastblock should be correct")

	err = s.Manager.SetLastBlockNumber(context.Background(), chain, updatedCurrentBlock)
}

func (s *ManagerTestSuite) TestManagerLastIndex() {
	chain := &dynamic.Chain{
		UUID: mockChain.UUID,
	}

	lastTxIndex, err := s.Manager.GetLastTxIndex(context.Background(), chain, 10)
	assert.NoError(s.T(), err, "GetLastTxIndex should not error")
	assert.Equal(s.T(), uint64(0), lastTxIndex, "LastTxIndex should be correct")

	err = s.Manager.SetLastTxIndex(context.Background(), chain, 10, 17)
	assert.NoError(s.T(), err, "SetLastTxIndex should not error")

	lastTxIndex, err = s.Manager.GetLastTxIndex(context.Background(), chain, 10)
	assert.NoError(s.T(), err, "GetLastTxIndex should not error")
	assert.Equal(s.T(), uint64(17), lastTxIndex, "LastTxIndex should be correct")

	lastTxIndex, err = s.Manager.GetLastTxIndex(context.Background(), chain, 11)
	assert.NoError(s.T(), err, "GetLastTxIndex should not error")
	assert.Equal(s.T(), uint64(0), lastTxIndex, "LastTxIndex should be correct")
}

func TestRegistry(t *testing.T) {
	suite.Run(t, new(ManagerTestSuite))
}
