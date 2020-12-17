package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
)

func TestAccountsParser(t *testing.T) {
	idenEntity := testutils.FakeAccount()
	idenModel := NewAccountModelFromEntities(idenEntity)
	finalIdenEntity := NewAccountEntityFromModels(idenModel)

	assert.Equal(t, idenEntity, finalIdenEntity)
}
