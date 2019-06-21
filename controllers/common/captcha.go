package common

import (
	"apcc_wallet/services/com"
	"apcc_wallet/utils"

	"github.com/mojocn/base64Captcha"

	"github.com/gin-gonic/gin"
)

var configD = base64Captcha.ConfigDigit{
	Height:     100,
	Width:      300,
	MaxSkew:    0.7,
	DotCount:   80,
	CaptchaLen: 6,
}

const (
	captchaVerifySuccess = "验证码校验成功"
	captchaVerifyFail    = "验证码校验失败"
	captchaGetSuccess    = "获取图片验证码成功"
)

var smsService com.SMSService

//TermController 期限结构控制器
type CaptchaController struct{}

func (CaptchaController) getCaption(c *gin.Context) {
	phone, hasPhone := c.GetQuery("phone")
	// device, hasDevice := c.GetQuery("device")
	// ip := c.Request.RemoteAddr
	var base64blob, captchaId string
	var captcaInterfaceInstance base64Captcha.CaptchaInterface
	if hasPhone {
		captchaId, captcaInterfaceInstance = base64Captcha.GenerateCaptcha(phone, configD)
		base64blob = base64Captcha.CaptchaWriteToBase64Encoding(captcaInterfaceInstance)
		utils.Write(c, true, captchaGetSuccess, map[string]interface{}{"img": base64blob, "captchaId": captchaId})
	}

}

func (CaptchaController) verificationCaption(c *gin.Context) {
	phone, hasPhone := c.GetQuery("phone")
	value, hasDevice := c.GetQuery("value")
	// ip := c.Request.RemoteAddr

	if hasPhone && hasDevice {
		verifyResult := base64Captcha.VerifyCaptcha(phone, value)
		if verifyResult {
			if err := smsService.SendSMS(phone); err == nil {
				utils.Write(c, true, captchaVerifySuccess, verifyResult)
			} else {
				utils.Write(c, false, err.Error(), nil)
			}

			return
		}
		utils.Write(c, false, captchaVerifyFail, verifyResult)

	}

}

func (this CaptchaController) Controller(c *gin.Context) {
	method := c.Request.Method
	switch method {
	case "GET":
		this.getCaption(c)
	case "POST":
		this.verificationCaption(c)
	}

}
