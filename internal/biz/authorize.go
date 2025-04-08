package biz

import (
	"coauth/internal/entities"
	"coauth/internal/model"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthorizeRepo interface {
	FindClientByID(ctx context.Context, clientID string) (*entities.App, error)
	FindAuthCode(ctx context.Context, code string) (*entities.AuthCode, error)
	SaveAuthCode(ctx context.Context, authCode *entities.AuthCode) error
	SaveToken(ctx context.Context, token *entities.Token) error
	FindToken(ctx context.Context, accessToken string) (*entities.Token, error)
	SaveClient(ctx context.Context, client *entities.App) error
}

type AuthorizeUseCase struct {
	repo AuthorizeRepo
	log  *zap.Logger
}

func NewAuthorizeUseCase(repo AuthorizeRepo, log *zap.Logger) *AuthorizeUseCase {
	return &AuthorizeUseCase{repo: repo, log: log}
}

func (uc *AuthorizeUseCase) Authorize(ctx context.Context, req *model.AuthorizeRequest) (*model.AuthorizeReply, error) {
	// 检查客户端是否存在
	client, err := uc.repo.FindClientByID(ctx, req.ClientID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "client not found")
	}

	// 生成授权码
	code := uuid.New().String()
	userID := ""
	expiresAt := time.Now().Add(time.Minute * 10)

	if req.RedirectURI == "" {
		req.RedirectURI = client.RedirectUrl
	}

	// save auth code
	authCode := &entities.AuthCode{
		Code:        code,
		ClientID:    req.ClientID,
		UserID:      userID,
		RedirectURL: req.RedirectURI,
		ExpiresAt:   expiresAt,
	}
	err = uc.repo.SaveAuthCode(ctx, authCode)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save auth code")
	}

	// 生成重定向URL
	redirectUrl := fmt.Sprintf("%s?code=%s&state=%s", authCode.RedirectURL, code, req.State)

	return &model.AuthorizeReply{
		Code:        code,
		State:       req.State,
		RedirectURI: redirectUrl,
	}, nil
}

func (uc *AuthorizeUseCase) Token(ctx context.Context, req *model.TokenRequest) (*model.TokenReply, error) {
	if req.GrantType == "authorization_code" {
		if req.Code == "" {
			return nil, status.Errorf(codes.InvalidArgument, "authorization code is required")
		}
	} else if req.GrantType == "refresh_token" {
		if req.RefreshToken == "" {
			return nil, status.Errorf(codes.InvalidArgument, "refresh token is required")
		}
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "invalid grant type")
	}

	if req.DataType == "jsonp" {
		if req.Callback == "" {
			return nil, status.Errorf(codes.InvalidArgument, "callback is required")
		}
	}

	// 检查客户端是否存在
	client, err := uc.repo.FindClientByID(ctx, req.ClientID)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid client id")
	}
	// 检查客户端密钥是否正确
	if client.Secret != req.ClientSecret {
		return nil, status.Errorf(codes.Unauthenticated, "invalid client credentials")
	}

	// 检查授权码是否存在
	authCode, err := uc.repo.FindAuthCode(ctx, req.Code)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization code")
	}
	// 检查授权码是否过期
	if authCode.ExpiresAt.Before(time.Now()) {
		return nil, status.Errorf(codes.Unauthenticated, "authorization code expired")
	}

	accessToken := uuid.New().String()
	refreshToken := uuid.New().String()
	tokenType := "Bearer"
	expiresIn := int64(3600)
	scopes := []string{}

	// 保存token
	token := &entities.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ClientID:     authCode.ClientID,
		UserID:       authCode.UserID,
		ExpiresAt:    time.Now().Add(time.Second * time.Duration(expiresIn)),
	}
	err = uc.repo.SaveToken(ctx, token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save token")
	}

	// TODO: 返回jsonp格式

	return &model.TokenReply{
		AccessToken:  accessToken,
		TokenType:    tokenType,
		ExpiresIn:    expiresIn,
		RefreshToken: refreshToken,
		Scope:        strings.Join(scopes, ","),
	}, nil
}

func (uc *AuthorizeUseCase) AddClient(ctx context.Context, req *model.AddClientRequest) (*model.AddClientReply, error) {
	// 生成客户端ID和密钥
	clientID := uuid.New().String()
	clientSecret := uuid.New().String()

	// 创建客户端对象
	client := &entities.App{
		ID:          clientID,
		Secret:      clientSecret,
		Name:        req.ClientName,
		RedirectUrl: req.CallbackURL,
		HomePage:    req.HomePage,
		Intro:       req.ClientSummary,
		Scopes:      strings.Join(req.Scopes, ","),
		UserId:      req.UserID,
	}

	// 保存客户端信息
	err := uc.repo.SaveClient(ctx, client)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to save client")
	}

	return &model.AddClientReply{
		ClientID: clientID,
	}, nil
}
