package middleware

import (
	"context"
	"net/http"
)

//go:generate mockgen -source=builder.go -destination=mock/builder.go -package=mock

type Builder interface {
	Build(ctx context.Context, name string, configuration interface{}) (mid func(http.Handler) http.Handler, respModifier func(resp *http.Response) error, err error)
}
