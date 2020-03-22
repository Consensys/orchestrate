// +build unit

package provider

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	providerName := "provider-test"
	qualifiedName := QualifyName(providerName, "test")
	assert.Equal(t, "test@provider-test", qualifiedName, "QualifyName should be correct")

	assert.Equal(t, providerName, GetName(qualifiedName), "GetName should be valid")

	ctx := WithName(context.Background(), qualifiedName)
	assert.Equal(t, providerName, NameFromContext(ctx), "Provider should have been attached to context")

	qualifiedName = QualifyNameFromContext(ctx, "test")
	assert.Equal(t, "test@provider-test", qualifiedName, "QualifyNameFromContext should be correct")
}
