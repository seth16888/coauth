package service

import (
	pb "coauth/api/v1"
	"coauth/internal/model"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *CoauthService) Captcha(context.Context, *pb.CaptchaRequest) (*pb.CaptchaReply, error) {
	key, value, err := s.captcha.GenerateCaptcha()
	if err != nil {
		return nil, err
	}
	return &pb.CaptchaReply{
		CaptchaKey:   key,
		CaptchaValue: value,
	}, nil
}

func (s *CoauthService) VerifyCaptcha(ctx context.Context, req *pb.VerifyCaptchaRequest) (
	*pb.VerifyCaptchaReply, error) {
	var params model.VerifyCaptchaReq
	params.CaptchaKey = req.CaptchaKey
	params.CaptchaCode = req.CaptchaCode
	params.Clear = req.Clear

	if err := s.validator.Validate(params); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if ok := s.captcha.VerifyCaptcha(
		params.CaptchaKey, params.CaptchaCode, params.Clear); !ok {
		return nil, status.Errorf(10400, "验证码错误")
	}
	return &pb.VerifyCaptchaReply{
		Success: true,
	}, nil
}
