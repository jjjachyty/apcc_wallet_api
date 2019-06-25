package commonCtr

import (
	"apcc_wallet_api/utils"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

//TermController 期限结构控制器
type SMSController struct{}

func (SMSController) verificationSMS(c *gin.Context) {
	var err = errors.New("短信验证码校验失败")
	phone, hasPhone := c.GetPostForm("phone")
	value, hasValue := c.GetPostForm("value")
	fmt.Println("verificationSMS", phone, value)

	if hasPhone && hasValue {
		err = smsService.VerificationSMS(phone, value)
	}
	utils.Response(c, err, nil)
}

func (this SMSController) Controller(c *gin.Context) {
	method := c.Request.Method
	switch method {
	case "POST":
		this.verificationSMS(c)
	}

}
