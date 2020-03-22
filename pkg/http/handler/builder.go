package handler

import (
	"context"
	"net/http"
)

//go:generate mockgen -destination=mock/handler.go -package=mock net/http Handler

//go:generate mockgen -source=builder.go -destination=mock/builder.go -package=mock

type Builder interface {
	Build(ctx context.Context, name string, configuration interface{}, respModifier func(*http.Response) error) (http.Handler, error)
}
