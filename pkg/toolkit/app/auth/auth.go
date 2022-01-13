package auth

import (
	"context"

	"github.com/consensys/orchestrate/pkg/toolkit/app/multitenancy"
)

//go:generate mockgen -source=auth.go -destination=mock/mock.go -package=mock

type Checker interface {
	Check(ctx context.Context) (*multitenancy.UserInfo, error)
}

type combinedChecker struct {
	checkers []Checker
}

func (a *combinedChecker) Check(ctx context.Context) (*multitenancy.UserInfo, error) {
	for _, checker := range a.checkers {
		userInfo, err := checker.Check(ctx)
		if err != nil {
			return nil, err
		}
		if userInfo != nil {
			return userInfo, nil
		}
	}

	return nil, nil
}

func NewCombineCheckers(checkers ...Checker) Checker {
	return &combinedChecker{checkers: checkers}
}
