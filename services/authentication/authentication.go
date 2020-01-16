package authentication

import "context"

type AuthFunc func(ctx context.Context) (context.Context, error)

type Auth interface {
	Check(ctx context.Context) (context.Context, error)
}

type combinedAuth struct {
	auths []Auth
}

func (a *combinedAuth) Check(ctx context.Context) (context.Context, error) {
	var err error
	for _, auth := range a.auths {
		ctx, err = auth.Check(ctx)
		if err == nil {
			return ctx, nil
		}
	}
	return ctx, err
}

func CombineAuth(auths ...Auth) Auth {
	return &combinedAuth{auths: auths}
}
