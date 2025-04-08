package service

import (
	pb "coauth/api/v1"
	"coauth/internal/biz"
	"coauth/internal/model"
	"context"

	"go.uber.org/zap"
)

type CoauthService struct {
	pb.UnimplementedCoauthServer
	uc  *biz.AuthorizeUseCase
	log *zap.Logger
}

func NewCoauthService(uc *biz.AuthorizeUseCase, log *zap.Logger) *CoauthService {
	return &CoauthService{
		uc:  uc,
		log: log,
	}
}

func (s *CoauthService) Authorize(ctx context.Context, req *pb.AuthorizeRequest) (*pb.AuthorizeReply, error) {
	s.log.Debug("Authorize", zap.Any("req", req))

	params := &model.AuthorizeRequest{
		ClientID:     req.ClientId,
		RedirectURI:  req.RedirectUri,
		ResponseType: req.ResponseType,
		State:        req.State,
	}
	authCode, err := s.uc.Authorize(ctx, params)
	if err != nil {
		return nil, err
	}

	return &pb.AuthorizeReply{
		Code:        authCode.Code,
		State:       authCode.State,
		RedirectUri: authCode.RedirectURI,
	}, nil
}

func (s *CoauthService) Token(ctx context.Context, req *pb.TokenRequest) (*pb.TokenReply, error) {
	s.log.Debug("Token", zap.Any("req", req))
	params := &model.TokenRequest{
		ClientID:     req.ClientId,
		ClientSecret: req.ClientSecret,
		GrantType:    req.GrantType,
		Code:         req.Code,
		RedirectURI:  req.RedirectUri,
		DataType:     req.DataType,
		Callback:     req.Callback,
		RefreshToken: req.RefreshToken,
	}

	token, err := s.uc.Token(ctx, params)
	if err != nil {
		return nil, err
	}

	return &pb.TokenReply{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    token.ExpiresIn,
		TokenType:    token.TokenType,
	}, nil
}

func (s *CoauthService) AddClient(ctx context.Context, req *pb.AddAppRequest) (*pb.AddAppReply, error) {
	s.log.Debug("AddClient", zap.Any("req", req))
	params := &model.AddClientRequest{
		ClientName:    req.ClientName,
		HomePage:      req.HomePage,
		ClientSummary: req.ClientSummary,
		CallbackURL:   req.CallbackUrl,
		Scopes:        req.Scopes,
		UserID:        req.UserId,
	}

	app, err := s.uc.AddClient(ctx, params)
	if err != nil {
		return nil, err
	}

	return &pb.AddAppReply{
		ClientId: app.ClientID,
		Code:     200,
		Message:  "success",
	}, nil
}
