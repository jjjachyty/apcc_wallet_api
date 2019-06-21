package common

import (
	"apcc_wallet/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	smsVerifySuccess = "验证码校验成功"
	smsVerifyFail    = "验证码校验失败"
)

//TermController 期限结构控制器
type SMSController struct{}

func (SMSController) verificationSMS(c *gin.Context) {
	phone, hasPhone := c.GetQuery("phone")
	value, hasValue := c.GetQuery("value")
	fmt.Println("verificationSMS", phone, value)

	if hasPhone && hasValue {
		if smsService.VerificationSMS(phone, value) {
			utils.Write(c, true, smsVerifySuccess, true)
			return
		}
		utils.Write(c, false, smsVerifyFail, false)

	}

}

func (this SMSController) Controller(c *gin.Context) {
	method := c.Request.Method
	switch method {
	case "POST":
		this.verificationSMS(c)
	}

}
