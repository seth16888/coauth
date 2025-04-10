package service

import (
	pb "coauth/api/v1"
	"context"
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
