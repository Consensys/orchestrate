package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/v2/pkg/types/testutils"
)

func TestAccountsParser(t *testing.T) {
	account := testutils.FakeAccount()
	accountModel := NewAccountModelFromEntities(account)
	finalAccount := NewAccountEntityFromModels(accountModel)

	assert.Equal(t, account, finalAccount)
}
