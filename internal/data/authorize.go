package data

import (
	"context"
	"errors"
	"fmt"

	"github.com/seth16888/coauth/internal/biz"
	"github.com/seth16888/coauth/internal/entities"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthorizeData struct {
	DB  *gorm.DB
	log *zap.Logger
}

func NewAuthorizeData(db *gorm.DB, log *zap.Logger) biz.AuthorizeRepo {
	return &AuthorizeData{DB: db, log: log}
}

// 根据客户端ID获取客户端信息
func (r *AuthorizeData) FindClientByID(ctx context.Context, clientID string) (*entities.App, error) {
	var client entities.App
	if err := r.DB.WithContext(ctx).Where("id = ?", clientID).First(&client).Error; err != nil {
		r.log.Error("failed to find client", zap.Error(err))
		return nil, fmtDBError(err)
	}
	return &client, nil
}

// 根据授权码获取授权信息
func (r *AuthorizeData) FindAuthCode(ctx context.Context, code string) (*entities.AuthCode, error) {
	var authCode entities.AuthCode
	if err := r.DB.WithContext(ctx).Where("code = ?", code).First(&authCode).Error; err != nil {
		r.log.Error("failed to find auth code", zap.Error(err))
		return nil, fmtDBError(err)
	}
	return &authCode, nil
}

// 保存授权码
func (r *AuthorizeData) SaveAuthCode(ctx context.Context, authCode *entities.AuthCode) error {
	err := r.DB.WithContext(ctx).Create(authCode).Error
	if err != nil {
		r.log.Error("failed to save auth code", zap.Error(err))
		return fmtDBError(err)
	}
	return nil
}

// 保存令牌
func (r *AuthorizeData) SaveToken(ctx context.Context, token *entities.Token) error {
	err := r.DB.WithContext(ctx).Create(token).Error
	if err != nil {
		r.log.Error("failed to save token", zap.Error(err))
		return fmtDBError(err)
	}
	return nil
}

// 根据访问令牌获取令牌信息
func (r *AuthorizeData) FindToken(ctx context.Context, accessToken string) (*entities.Token, error) {
	var token entities.Token
	if err := r.DB.WithContext(ctx).Where("access_token = ?", accessToken).First(&token).Error; err != nil {
		r.log.Error("failed to find token", zap.Error(err))
		return nil, fmtDBError(err)
	}
	return &token, nil
}

func (r *AuthorizeData) SaveClient(ctx context.Context, client *entities.App) error {
	err := r.DB.WithContext(ctx).Create(client).Error
	if err != nil {
		r.log.Error("failed to save client", zap.Error(err))
		return fmtDBError(err)
	}
	return nil
}

func fmtDBError(err error) error {
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("database occur error: %s", err.Error())
	}
	return fmt.Errorf("data not found")
}
