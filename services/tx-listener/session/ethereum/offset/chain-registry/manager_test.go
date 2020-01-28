package chainregistry

import (
	"context"
	"fmt"
	"testing"

	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/chain-registry/store/types"

	"github.com/stretchr/testify/suite"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/tx-listener/session/ethereum/offset/testutils"
)

var MockChainsSlice = []*types.Chain{
	{
		UUID:                    "test-chain",
		Name:                    "test",
		TenantID:                "test",
		URLs:                    []string{"test"},
		ListenerDepth:           &(&struct{ x uint64 }{0}).x,
		ListenerBlockPosition:   &(&struct{ x int64 }{0}).x,
		ListenerFromBlock:       &(&struct{ x int64 }{0}).x,
		ListenerBackOffDuration: &(&struct{ x string }{"0s"}).x,
	},
	{
		UUID:                    "test-chain1",
		Name:                    "test1",
		TenantID:                "test1",
		URLs:                    []string{"test1"},
		ListenerDepth:           &(&struct{ x uint64 }{1}).x,
		ListenerBlockPosition:   &(&struct{ x int64 }{1}).x,
		ListenerFromBlock:       &(&struct{ x int64 }{1}).x,
		ListenerBackOffDuration: &(&struct{ x string }{"1s"}).x,
	},
}

var MockChainsMap = map[string]*types.Chain{
	"test-chain":  MockChainsSlice[0],
	"test-chain1": MockChainsSlice[1],
}

type Mock struct{}

func (c *Mock) GetChainByUUID(_ context.Context, chainUUID string) (*types.Chain, error) {
	if _, ok := MockChainsMap[chainUUID]; !ok {
		return nil, fmt.Errorf("test")
	}
	return MockChainsMap[chainUUID], nil
}

func (c *Mock) GetChainByTenantAndName(_ context.Context, _, _ string) (*types.Chain, error) {
	return nil, nil
}

func (c *Mock) GetChainByTenantAndUUID(_ context.Context, _, _ string) (*types.Chain, error) {
	return nil, nil
}

func (c *Mock) GetChains(_ context.Context) ([]*types.Chain, error) {
	return MockChainsSlice, nil
}

func (c *Mock) UpdateBlockPosition(_ context.Context, chainUUID string, blockNumber int64) error {
	if _, ok := MockChainsMap[chainUUID]; !ok {
		return fmt.Errorf("test")
	}
	MockChainsMap[chainUUID].ListenerBlockPosition = &(&struct{ x int64 }{blockNumber}).x
	return nil
}

type ManagerTestSuite struct {
	testutils.OffsetManagerTestSuite
}

func (s *ManagerTestSuite) SetupTest() {
	s.Manager = NewManager(&Mock{})
}

func TestRegistry(t *testing.T) {
	s := new(ManagerTestSuite)
	suite.Run(t, s)
}
