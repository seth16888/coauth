package data

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type GORMRepository struct {
	DB *gorm.DB
}

func NewGORMRepository(db *gorm.DB) *GORMRepository {
	return &GORMRepository{DB: db}
}

// 根据客户端ID获取客户端信息
func (r *GORMRepository) FindClientByID(ctx context.Context, clientID string) (*App, error) {
	var client App
	if err := r.DB.WithContext(ctx).Where("id = ?", clientID).First(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &client, nil
}

// 根据授权码获取授权信息
func (r *GORMRepository) FindAuthCode(ctx context.Context, code string) (*AuthCode, error) {
	var authCode AuthCode
	if err := r.DB.WithContext(ctx).Where("code = ?", code).First(&authCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &authCode, nil
}

// 保存授权码
func (r *GORMRepository) SaveAuthCode(ctx context.Context, authCode *AuthCode) error {
	return r.DB.WithContext(ctx).Create(authCode).Error
}

// 保存令牌
func (r *GORMRepository) SaveToken(ctx context.Context, token *Token) error {
	return r.DB.WithContext(ctx).Create(token).Error
}

// 根据访问令牌获取令牌信息
func (r *GORMRepository) FindToken(ctx context.Context, accessToken string) (*Token, error) {
	var token Token
	if err := r.DB.WithContext(ctx).Where("access_token = ?", accessToken).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func (r *GORMRepository) SaveClient(ctx context.Context, client *App) error {
	return r.DB.WithContext(ctx).Create(client).Error
}
