package controllers

import (
	"github.com/mojocn/base64Captcha"
)

var store = base64Captcha.DefaultMemStore

func GenerateCaptcha() (string, string, error) {
	driver := base64Captcha.NewDriverDigit(80, 240, 5, 0.7, 80) // 图形验证码配置
	captcha := base64Captcha.NewCaptcha(driver, store)

	id, b64s, _, err := captcha.Generate()
	if err != nil {
		return "", "", err
	}
	return id, b64s, nil
}

func VerifyCaptcha(captchaId string, value string) (bool, error) {
	if store.Verify(captchaId, value, true) {
		return true, nil
	} else {
		return false, nil
	}
}
