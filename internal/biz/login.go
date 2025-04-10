package biz

import (
	"coauth/internal/entities"
	"coauth/internal/model"
	"coauth/pkg/encrypt"
	"coauth/pkg/helpers"
	"coauth/pkg/jwt"
	"context"
	"fmt"

	goJWT "github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type LoginRepo interface {
	FindByName(ctx context.Context, name string) (*entities.User, error)
	UpdateLoginInfo(ctx context.Context, uid string, lastLoginTime int64, lastLoginIp string) error
}

type LoginUsecase struct {
	repo LoginRepo
	log  *zap.Logger
	jwt  *jwt.JWTService
	bk   *TokenBlacklistUsecase
}

func NewLoginUsecase(repo LoginRepo, log *zap.Logger,
	jwt *jwt.JWTService, bk *TokenBlacklistUsecase) *LoginUsecase {
	return &LoginUsecase{repo: repo, log: log, jwt: jwt, bk: bk}
}

func (uc *LoginUsecase) PasswordLogin(ctx context.Context,
	req *model.PasswordLoginReq) (*model.PasswordLoginRes, error) {
	uc.log.Debug("PasswordLogin", zap.Any("req", req))

	user, err := uc.repo.FindByName(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if user.Status == nil || *user.Status == 0 {
		return nil, fmt.Errorf("user is disabled")
	}

	if encrypt.EncryptPassword(req.Password, user.Salt) != user.Password {
		return nil, fmt.Errorf("username or password error")
	}

	// 生成access token
	now := helpers.TimenowInTimezone(uc.jwt.TimeLocation)
	claims := &jwt.JWTCustomClaims{
		UserID:       user.ID,
		DepartmentID: 0,
		Authorities:  []string{"ROLE_Admin"},
		DataScope:    1,
		RegisteredClaims: goJWT.RegisteredClaims{
			Subject:   user.Username,
			Issuer:    uc.jwt.Issuer,                                    // 签发人
			ExpiresAt: goJWT.NewNumericDate(now.Add(uc.jwt.ExpireTime)), // 过期时间
			IssuedAt:  goJWT.NewNumericDate(now),                        // 签发时间
			NotBefore: goJWT.NewNumericDate(now),                        // 生效时间
			ID:        helpers.NewUUID(),                                // 编号
		},
	}
	ak, exp, err := uc.jwt.IssueToken(claims)
	if err != nil {
		return nil, err
	}

	return &model.PasswordLoginRes{
		AccessToken:  ak,
		RefreshToken: ak,
		TokenType:    "Bearer",
		ExpiresIn:    exp,
	}, nil
}

// Logout
func (uc *LoginUsecase) Logout(accessToken string) error {
	// 解析token
	claims, err := uc.jwt.ParseToken(accessToken)
	if err != nil {
		return err
	}
	uc.log.Debug("Logout", zap.String("user", claims.UserID))
	// token 放入黑名单
	if err := uc.bk.Push(context.Background(), claims.ID); err != nil {
		return err
	}
	return nil
}

// RefreshToken processes the refresh token logic.
func (uc *LoginUsecase) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	// Generate a new access token
	newAccessToken, err := uc.jwt.RefreshToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}
	uc.log.Debug("RefreshToken", zap.String("newAccessToken", newAccessToken))
	return newAccessToken, nil
}
