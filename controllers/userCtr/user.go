package userCtr

import (
	"apcc_wallet_api/middlewares/jwt"
	"apcc_wallet_api/models/userMod"
	"apcc_wallet_api/services/commonSrv"
	"apcc_wallet_api/services/userSrv"
	"apcc_wallet_api/utils"
	"errors"
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
)

var userService userSrv.UserService
var smsService commonSrv.SMSService

//TermController 期限结构控制器
type UserController struct{}

func (UserController) Register(c *gin.Context) {
	var err error
	var user = new(userMod.User)
	if err = c.BindJSON(user); err == nil {
		if VerifyMobileFormat(user.Phone) && VerifyPasswdFormat(user.Password) {
			err = userService.Register(user)
		}
	}
	utils.Response(c, err, nil)

}

//LoginWithPW 用户名密码登录
func (UserController) LoginWithPW(c *gin.Context) {
	var err error
	var token string
	var user = new(userMod.User)
	if err = c.BindJSON(user); err == nil {
		if VerifyMobileFormat(user.Phone) && VerifyPasswdFormat(user.Password) {
			err = userService.Login(user)
			if user.NickName == "" {
				err = errors.New("用户名或密码错误")
			} else {
				token, err = jwt.GenerateToken(*user)
			}

		}
	}
	utils.Response(c, err, map[string]interface{}{"User": user, "Token": token})
}

//LoginWithSMS 短信验证码登录
func (UserController) LoginWithSMS(c *gin.Context) {
	var err error
	var user = new(userMod.User)
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

func (UserController) PayPassword(c *gin.Context) {
	var err = errors.New("密码不能为空")
	orgPassword := c.PostForm("orgPassword")
	password := c.PostForm("password")
	fmt.Println("password=", password)
	claims := jwt.GetClaims(c)
	var newToken string
	if claims.HasPayPasswd {
		if orgPassword != "" && password != "" {

			//新增

		} else {
			err = errors.New("修改密码原密码和密码不能为空")
		}
	} else {
		if password != "" {
			if err = userService.Update(&userMod.User{Phone: claims.Phone, PayPasswd: utils.GetMD5(password)}); err == nil {
				claims.HasPayPasswd = true
				newToken, err = jwt.NewJWT().CreateToken(*claims)
			}
		}

	}

	utils.Response(c, err, gin.H{"Token": newToken})
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
