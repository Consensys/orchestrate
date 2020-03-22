package auth

import (
	"context"
)

//go:generate mockgen -source=auth.go -destination=mock/mock.go -package=mock

type Checker interface {
	Check(ctx context.Context) (context.Context, error)
}

type combinedChecker struct {
	checkers []Checker
}

func (a *combinedChecker) Check(ctx context.Context) (context.Context, error) {
	var err error
	for _, checker := range a.checkers {
		ctx, err = checker.Check(ctx)
		if err == nil {
			return ctx, nil
		}
	}
	return ctx, err
}

func CombineCheckers(checkers ...Checker) Checker {
	return &combinedChecker{checkers: checkers}
}
