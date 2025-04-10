package biz

import (
	"coauth/pkg/captcha"
	"image/color"

	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

type CaptchaUsecase struct {
	store captcha.Store
	log   *zap.Logger
	text  *captcha.TextCaptcha
}

func NewCaptchaUsecase(store captcha.Store, log *zap.Logger) *CaptchaUsecase {
	// driver
	h := 50
	w := 150
	noise := 2
	showline := 2
	length := 4
	source := "234567890abcdefghjkmnpqrstuvwxyz"
	// 字体
	//  RitaSmith.ttf | wqy-microhei.ttc | Flim-Flam.ttf | DENNEthree-dee.ttf
	//  DeborahFancyDress.ttf | Comismsh.ttf | chromohv.ttf | ApothecaryFont.ttf
	//  actionj.ttf | 3Dumb.ttf
	fonts := []string{"DENNEthree-dee.ttf"}
	bgColor := &color.RGBA{R: 240, G: 240, B: 246, A: 246}

	dri := base64Captcha.NewDriverString(h, w, noise, showline, length, source, bgColor, nil, fonts)

	c := &CaptchaUsecase{}
	c.text = captcha.NewTextCaptcha(store, dri)
	c.log = log
	c.store = store
	return c
}

func (uc *CaptchaUsecase) GenerateCaptcha() (key string, value string, err error) {
	key, value, _, err = uc.text.Make()
	if err != nil {
		uc.log.Error("GenerateCaptcha", zap.Error(err))
		return "", "", err
	}
	return key, value, nil
}

func (uc *CaptchaUsecase) VerifyCaptcha(key, answer string, clear bool) bool {
	// TODO: 测试代码，需要删除
	if answer == "1234" {
		return true
	}
	return uc.text.Verify(key, answer, clear)
}
