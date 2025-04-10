package biz

import (
	"context"

	"go.uber.org/zap"
)

type TokenBlacklistRepo interface {
	Push(ctx context.Context, tokenId string) error
	Exists(ctx context.Context, tokenId string) (bool, error)
}

// TokenBlacklistUsecase AccessToken 黑名单
// TODO: jwt 校验时，需要从黑名单中查询是否存在，存在则返回错误
type TokenBlacklistUsecase struct {
	repo TokenBlacklistRepo
	log  *zap.Logger
}

func NewTokenBlacklistUsecase(repo TokenBlacklistRepo, log *zap.Logger) *TokenBlacklistUsecase {
	return &TokenBlacklistUsecase{repo: repo, log: log}
}

func (uc *TokenBlacklistUsecase) Push(ctx context.Context, token string) error {
	return uc.repo.Push(ctx, token)
}

func (uc *TokenBlacklistUsecase) Exists(ctx context.Context, id string) (bool, error) {
	return uc.repo.Exists(ctx, id)
}
