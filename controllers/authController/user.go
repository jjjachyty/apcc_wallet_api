package authController

import (
	"apcc_wallet_api/models/authModel"
	"apcc_wallet_api/services/authService"
	"apcc_wallet_api/services/commonService"
	"apcc_wallet_api/utils"
	"errors"
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
)

var userService authService.UserService
var smsService commonService.SMSService

//TermController 期限结构控制器
type RegisterController struct{}

func (RegisterController) Register(c *gin.Context) {
	var err error
	var user = new(authModel.User)
	if err = c.BindJSON(user); err == nil {
		if VerifyMobileFormat(user.Phone) && VerifyPasswdFormat(user.Password) {
			err = userService.Register(user)
		}
	}
	utils.Response(c, err, nil)

}

//LoginWithPW 用户名密码登录
func (RegisterController) LoginWithPW(c *gin.Context) {
	var err error

	var user = new(authModel.User)
	if err = c.BindJSON(user); err == nil {
		if VerifyMobileFormat(user.Phone) && VerifyPasswdFormat(user.Password) {

			err = userService.Login(user)
			if err == nil && user.NickName == "" {
				err = errors.New("用户名或密码错误")
			}

		}
	}
	fmt.Println(user)
	utils.Response(c, err, nil)
}

//LoginWithSMS 短信验证码登录
func (RegisterController) LoginWithSMS(c *gin.Context) {
	var err error
	var user = new(authModel.User)
	phone := c.PostForm("phone")
	sms := c.PostForm("sms")
	fmt.Println(phone, sms)
	if VerifyMobileFormat(phone) {
		if err = smsService.VerificationSMS(phone, sms); err == nil {
			if err = userService.Get(user); err == nil {

			}
		}
	}
	fmt.Println(user)
	utils.Response(c, err, user)

}

func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func VerifyPasswdFormat(mobileNum string) bool {
	regular := "^[a-z0-9_-]{6,16}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}
