package data

import (
	"coauth/internal/biz"
	"coauth/internal/entities"
	"context"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type LoginData struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewLoginData(db *gorm.DB, logger *zap.Logger) biz.LoginRepo {
	return &LoginData{
		db:     db,
		logger: logger,
	}
}

// FindByName find user by name
func (a *LoginData) FindByName(ctx context.Context, name string) (*entities.User, error) {
	user := &entities.User{}
	query := a.db.Model(&entities.User{}).
		Where("username =?", name)
	if err := query.First(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateLoginInfo update user login info
func (a *LoginData) UpdateLoginInfo(ctx context.Context, uid string, lastActive int64, lastIp string) error {
	return a.db.Model(&entities.User{}).Where("id = ?", uid).Updates(
		map[string]any{"last_login_at": lastActive, "last_login_ip": lastIp},
	).Error
}
