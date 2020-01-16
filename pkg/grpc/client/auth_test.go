package grpcclient

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	grpcserver "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/pkg/server/grpc"
	authutils "gitlab.com/ConsenSys/client/fr/core-stack/orchestrate.git/services/authentication/utils"
)

func TestPerRPCCredentials(t *testing.T) {
	cred := &PerRPCCredentials{}
	h, _ := cred.GetRequestMetadata(context.Background(), "")
	assert.Len(t, h, 0, "Header length should be valid")

	ctx := authutils.WithAuthorization(context.Background(), "test-auth")
	h, _ = cred.GetRequestMetadata(ctx, "")
	assert.Len(t, h, 1, "Header length should be valid")
	assert.Equal(t, "test-auth", h[grpcserver.AuthorizationHeader], "Header should be correct")
}
