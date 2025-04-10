package service

import (
	pb "coauth/api/v1"
	"coauth/internal/model"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoauthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	// 参数转换
	var params model.PasswordLoginReq
	params.Username = req.Username
	params.Password = req.Password
	params.CaptchaKey = req.CaptchaKey
	params.CaptchaCode = req.CaptchaCode
	// 参数校验
	if err := s.validator.Validate(params); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// 登录
	res, err := s.login.PasswordLogin(ctx, &params)
	if err != nil {
		return nil, status.Errorf(10400, "invalid username or password")
	}

	return &pb.LoginReply{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		TokenType:    res.TokenType,
		ExpiresIn:    res.ExpiresIn,
	}, nil
}
