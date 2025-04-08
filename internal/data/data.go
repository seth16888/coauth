package data

import (
	"coauth/internal/biz"
	"coauth/internal/entities"
	"context"
	"errors"

	"gorm.io/gorm"
)

type AuthorizeData struct {
	DB *gorm.DB
}

func NewAuthorizeData(db *gorm.DB) biz.AuthorizeRepo {
	return &AuthorizeData{DB: db}
}

// 根据客户端ID获取客户端信息
func (r *AuthorizeData) FindClientByID(ctx context.Context, clientID string) (*entities.App, error) {
	var client entities.App
	if err := r.DB.WithContext(ctx).Where("id = ?", clientID).First(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &client, nil
}

// 根据授权码获取授权信息
func (r *AuthorizeData) FindAuthCode(ctx context.Context, code string) (*entities.AuthCode, error) {
	var authCode entities.AuthCode
	if err := r.DB.WithContext(ctx).Where("code = ?", code).First(&authCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &authCode, nil
}

// 保存授权码
func (r *AuthorizeData) SaveAuthCode(ctx context.Context, authCode *entities.AuthCode) error {
	return r.DB.WithContext(ctx).Create(authCode).Error
}

// 保存令牌
func (r *AuthorizeData) SaveToken(ctx context.Context, token *entities.Token) error {
	return r.DB.WithContext(ctx).Create(token).Error
}

// 根据访问令牌获取令牌信息
func (r *AuthorizeData) FindToken(ctx context.Context, accessToken string) (*entities.Token, error) {
	var token entities.Token
	if err := r.DB.WithContext(ctx).Where("access_token = ?", accessToken).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func (r *AuthorizeData) SaveClient(ctx context.Context, client *entities.App) error {
	return r.DB.WithContext(ctx).Create(client).Error
}
