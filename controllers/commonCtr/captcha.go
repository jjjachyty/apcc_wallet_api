package commonCtr

import (
	"apcc_wallet_api/services/commonSrv"
	"apcc_wallet_api/utils"
	"errors"

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

var smsService commonSrv.SMSService

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
		utils.Response(c, nil, map[string]interface{}{"img": base64blob, "captchaId": captchaId})
	}

}

func (CaptchaController) verificationCaption(c *gin.Context) {
	var err error
	phone := c.PostForm("phone")
	value := c.PostForm("value")
	// ip := c.Request.RemoteAddr

	if phone != "" && value != "" {
		verifyResult := base64Captcha.VerifyCaptchaAndIsClear(phone, value, true)
		if verifyResult {
			err = smsService.SendSMS(phone)

		} else {
			err = errors.New("校验验证码错误")
		}
	}
	utils.Response(c, err, nil)

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
