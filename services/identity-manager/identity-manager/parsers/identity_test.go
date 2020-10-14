package parsers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/types/testutils"
)

func TestIdentityParser(t *testing.T) {
	idenEntity := testutils.FakeIdentity()
	idenModel := NewIdentityModelFromEntities(idenEntity)
	finalIdenEntity := NewIdentityEntityFromModels(idenModel)

	assert.Equal(t, idenEntity, finalIdenEntity)
}
