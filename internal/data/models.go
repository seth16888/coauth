package data

import (
	"time"

	"gorm.io/gorm"
)

type App struct {
	ID          string         `gorm:"column:id;primaryKey;size:36"`
	Secret      string         `gorm:"column:secret;size:64;not null"`
	Name        string         `gorm:"column:name;size:32;unique;not null"`
	RedirectUrl string         `gorm:"column:redirect_url;size:256;not null"`
	HomePage    string         `gorm:"column:home_page;size:256;not null"`
	Intro       string         `gorm:"column:intro;size:64"`
	Scopes      string         `gorm:"column:scopes;size:256"`
	UserId      string         `gorm:"column:user_id;size:36;not null"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 自定义表名
func (App) TableName() string {
	return "app"
}

type User struct {
	ID          string         `gorm:"column:id;primaryKey;size:36"`
	Username    string         `gorm:"column:username;size:32;unique;not null"`
	Password    string         `gorm:"column:password;size:64;not null"`
	Salt        string         `gorm:"column:salt;size:16;not null"`
	Nickname    string         `gorm:"column:nickname;size:32"`
	Avatar      string         `gorm:"column:avatar;size:256"`
	Email       string         `gorm:"column:email;size:128"`
	Phone       string         `gorm:"column:phone;size:20"`
	Status      *int           `gorm:"column:status;not null"`
	LastLoginAt int64          `gorm:"column:last_login_at"`
	LastLoginIp string         `gorm:"column:last_login_ip;size:16"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 自定义表名
func (User) TableName() string {
	return "user"
}

type AuthCode struct {
	ID          int64          `gorm:"column:id;primaryKey;autoIncrement"`
	Code        string         `gorm:"column:code;size:36"`
	ClientID    string         `gorm:"column:client_id;size:36;not null"`
	UserID      string         `gorm:"column:user_id;size:36;not null"`
	RedirectURL string         `gorm:"column:redirect_url;size:256;not null"`
	ExpiresAt   time.Time      `gorm:"column:expires_at;not null"`
	CreatedAt   time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 自定义表名
func (AuthCode) TableName() string {
	return "auth_code"
}

type Token struct {
	ID           int64          `gorm:"column:id;primaryKey;autoIncrement"`
	AccessToken  string         `gorm:"column:access_token;size:36"`
	RefreshToken string         `gorm:"column:refresh_token;size:36"`
	ClientID     string         `gorm:"column:client_id;size:36;not null"`
	UserID       string         `gorm:"column:user_id;size:36;not null"`
	ExpiresAt    time.Time      `gorm:"column:expires_at;not null"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

// TableName 自定义表名
func (Token) TableName() string {
	return "token"
}
