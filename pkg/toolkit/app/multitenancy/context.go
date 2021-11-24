package multitenancy

import (
	"context"
)

type multitenancyCtxKey string

const (
	UserInfoKey multitenancyCtxKey = "user_info"
)

func WithUserInfo(ctx context.Context, userInfo *UserInfo) context.Context {
	return context.WithValue(ctx, UserInfoKey, userInfo)
}

func UserInfoValue(ctx context.Context) *UserInfo {
	userInfo, ok := ctx.Value(UserInfoKey).(*UserInfo)
	if !ok {
		return nil
	}
	return userInfo
}
