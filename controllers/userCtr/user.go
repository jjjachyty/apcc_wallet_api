package userCtr

import (
	"apcc_wallet_api/middlewares/jwt"
	"apcc_wallet_api/models/assetMod"
	"apcc_wallet_api/models/userMod"
	"apcc_wallet_api/services/assetSrv"
	"apcc_wallet_api/services/commonSrv"
	"apcc_wallet_api/services/userSrv"
	"apcc_wallet_api/services/walletSrv"
	"apcc_wallet_api/utils"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"regexp"

	"github.com/gin-gonic/gin"
)

var userService userSrv.UserService
var smsService commonSrv.SMSService
var assetsService assetSrv.AssetService

//TermController 期限结构控制器
type UserController struct{}

func (UserController) Register(c *gin.Context) {
	var err error
	var user = new(userMod.User)
	var assets []assetMod.Asset
	if err = c.BindJSON(user); err == nil {

		if VerifyMobileFormat(user.Phone) && VerifyPasswdFormat(user.Password) {
			session := utils.OpenSession()
			defer session.Close()
			user.UUID = utils.GetUUID()
			user.State = utils.STATE_ENABLE
			if err = userService.Register(user); err == nil {
				if assets, err = walletSrv.GetAddress(user.UUID, uint32(user.AccountID)); err == nil {
					err = assetsService.Create(assets)
				}
			}

		}
	}
	fmt.Println(err)
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
			fmt.Println(err, "UUID", user.UUID)
			if user.UUID != "" {
				fmt.Println("user.PayPasswd", user.PayPasswd)
				if user.PayPasswd != "" {
					user.HasPayPasswd = true
				}

				token, err = jwt.GenerateToken(*user)
			} else {
				err = errors.New("用户名或密码错误")
			}

		} else {
			err = errors.New("用户名或密码格式错误")
		}
	}
	fmt.Println(err, user)
	utils.Response(c, err, map[string]interface{}{"User": gin.H{"IDCardAuth": user.IDCardAuth, "UUID": user.UUID, "Phone": user.Phone, "HasPayPasswd": user.HasPayPasswd, "NickName": user.NickName, "Avatar": user.Avatar}, "Token": token})
}

//LoginWithSMS 短信验证码登录
func (UserController) LoginWithSMS(c *gin.Context) {
	var err error
	var user = new(userMod.User)
	phone := c.PostForm("phone")
	sms := c.PostForm("sms")
	var token string
	fmt.Println(phone, sms)
	if VerifyMobileFormat(phone) {
		if err = smsService.VerificationSMS(phone, sms); err == nil {
			if err = userService.Get(user); err == nil {
				if user.UUID != "" {
					if user.PayPasswd != "" {
						user.HasPayPasswd = true
					}

					token, err = jwt.GenerateToken(*user)
				}
			}
		}
	}
	fmt.Println(user)
	utils.Response(c, err, map[string]interface{}{"User": gin.H{"IDCardAuth": user.IDCardAuth, "UUID": user.UUID, "Phone": user.Phone, "HasPayPasswd": user.HasPayPasswd, "NickName": user.NickName, "Avatar": user.Avatar}, "Token": token})

}

func (UserController) PayPassword(c *gin.Context) {
	var err = errors.New("密码不能为空")
	orgPassword := c.PostForm("orgPassword")
	password := c.PostForm("password")
	claims := jwt.GetClaims(c)
	var newToken string

	if claims.HasPayPasswd {
		if orgPassword != "" && password != "" {
			var user = new(userMod.User)
			user.Phone = claims.Phone
			user.PayPasswd = utils.GetMD5(orgPassword)
			if err = userService.Get(user); err == nil {
				fmt.Println(user.UUID)
				if user.UUID != "" {
					user.PayPasswd = utils.GetMD5(password)
					err = userService.Update(user)
				} else {
					err = errors.New("原密码错误")
				}
			}
			//新增

		} else {
			err = errors.New("修改密码原密码和密码不能为空")
		}
		utils.Response(c, err, true)
	} else {
		if password != "" {
			if err = userService.Update(&userMod.User{Phone: claims.Phone, PayPasswd: utils.GetMD5(password)}); err == nil {
				claims.HasPayPasswd = true
				newToken, err = jwt.NewJWT().CreateToken(*claims)
			}
		}
		utils.Response(c, err, gin.H{"Token": newToken})
	}

}
func (UserController) LoginPassword(c *gin.Context) {
	var err = errors.New("密码不能为空")
	password := c.PostForm("password")
	claims := jwt.GetClaims(c)

	if password != "" {
		var user = new(userMod.User)
		user.Phone = claims.Phone
		user.Password = utils.GetMD5(password)
		//新增
		err = userService.Update(user)
	}

	utils.Response(c, err, nil)

}

func (UserController) Profile(c *gin.Context) {
	var err error
	var fh *multipart.FileHeader
	var file multipart.File
	var imageBytes []byte

	nickName := c.PostForm("nickName")
	claims := jwt.GetClaims(c)
	if fh, _ = c.FormFile("avatar"); fh != nil {
		if file, err = fh.Open(); err == nil {
			defer file.Close()
			imageBytes, err = ioutil.ReadAll(file)
			err = commonSrv.UploadImage(claims.UUID, imageBytes)
		}
	}
	if err == nil && nickName != "" {
		err = userService.Update(&userMod.User{UUID: claims.UUID, Avatar: claims.UUID, NickName: nickName})
	}
	utils.Response(c, err, nil)
}

func (UserController) IDCard(c *gin.Context) {
	var err error
	card := new(userMod.IdCard)
	if err = c.BindJSON(card); err == nil {
		if card.Name != "" && card.IDCardNumber != "" {
			card.UserID = jwt.GetClaims(c).UUID
			err = userService.UpdateIDCard(card)

		} else {
			utils.Response(c, errors.New("参数错误"), nil)
		}
	}
	fmt.Println(err)
	utils.Response(c, err, nil)
}

func VerifyMobileFormat(mobileNum string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func VerifyPasswdFormat(passwd string) bool {
	regular := "^[0-9A-Za-z].{8,16}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(passwd)
}
