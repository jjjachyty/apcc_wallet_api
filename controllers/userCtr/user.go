package userCtr

import (
	"apcc_wallet_api/middlewares/jwt"
	"apcc_wallet_api/models/userMod"
	"apcc_wallet_api/services/assetSrv"
	"apcc_wallet_api/services/commonSrv"
	"apcc_wallet_api/services/userSrv"
	"apcc_wallet_api/utils"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
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
	if err = c.BindJSON(user); err == nil {

		if VerifyMobileFormat(user.Phone) && VerifyPasswdFormat(user.Password) {
			err = userService.Register(user)
		}

	}
	utils.Response(c, err, nil)

}

//CheckPhone 检查手机号是否已注册
func (UserController) CheckPhone(c *gin.Context) {
	var err error
	phone, hasPhone := c.GetQuery("phone")
	if hasPhone {
		user := new(userMod.User)
		user.Phone = phone
		if err = userService.Get(user); err == nil {
			if user.UUID != "" {
				err = errors.New("用户已存在,请登录")
			}
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
	utils.Response(c, err, map[string]interface{}{"User": gin.H{"IDCardAuth": user.IDCardAuth, "UUID": user.UUID, "Phone": user.Phone, "Introduce": user.Introduce, "HasPayPasswd": user.HasPayPasswd, "NickName": user.NickName, "Avatar": user.Avatar}, "Token": token})
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
			user.Phone = phone
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
	utils.Response(c, err, map[string]interface{}{"User": gin.H{"IDCardAuth": user.IDCardAuth, "UUID": user.UUID, "Phone": user.Phone, "Introduce": user.Introduce, "HasPayPasswd": user.HasPayPasswd, "NickName": user.NickName, "Avatar": user.Avatar}, "Token": token})

}

func (UserController) PayPassword(c *gin.Context) {
	var err = errors.New("密码不能为空")

	password := c.PostForm("payPassword")
	claims := jwt.GetClaims(c)

	// if claims.HasPayPasswd {
	if password != "" {
		var user = new(userMod.User)
		user.UUID = claims.UUID
		user.PayPasswd = utils.GetMD5(password)
		err = userService.Update(user)
	}

	//新增

	utils.Response(c, err, true)
	// } else {
	// if password != "" {
	// 	if err = userService.Update(&userMod.User{Phone: claims.Phone, PayPasswd: utils.GetMD5(password)}); err == nil {

	// 		newToken, err = jwt.NewJWT().CreateToken(*claims)
	// 	}
	// }
	// utils.Response(c, err, gin.H{"Token": newToken})
	// }

}
func (UserController) LoginPassword(c *gin.Context) {
	var err = errors.New("密码不能为空")
	password := c.PostForm("password")
	claims := jwt.GetClaims(c)

	if password != "" {
		var user = new(userMod.User)
		user.UUID = claims.UUID
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
	var user userMod.User
	nickName := c.PostForm("nickName")
	introduce := c.PostForm("introduce")
	claims := jwt.GetClaims(c)
	if fh, _ = c.FormFile("avatar"); fh != nil {
		if file, err = fh.Open(); err == nil {
			defer file.Close()
			imageBytes, err = ioutil.ReadAll(file)
			err = commonSrv.UploadImage(claims.UUID, imageBytes)
		}
	}
	if err == nil {
		user.UUID = claims.UUID
		user.Avatar = claims.UUID
		if nickName != "" {
			user.NickName = nickName
		}
		if introduce != "" {
			user.Introduce = introduce
		}
		err = userService.Update(&user)
	}
	utils.Response(c, err, nil)
}

func (UserController) IDCard(c *gin.Context) {
	var err error
	card := new(userMod.IdCard)
	if err = c.BindJSON(card); err == nil {
		if card.Name != "" && card.IDCardNumber != "" {
			card.UserID = jwt.GetClaims(c).UUID
			if err = userService.AddIDCard(card); err == nil {
				err = userService.Update(&userMod.User{UUID: card.UserID, IDCardAuth: 1})
			}

		} else {
			utils.Response(c, errors.New("参数错误"), nil)
		}
	}
	utils.Response(c, err, nil)
}

func (UserController) GetUserByPhone(c *gin.Context) {
	var err error
	var phoneBytes []byte
	var user = new(userMod.User)
	phone := c.Param("id")
	fmt.Println(c.Param("id"), c.Query("id"))
	if phone != "" {
		if len(phone) != 11 {
			if phoneBytes, err = hex.DecodeString(phone); err == nil {
				phone = big.NewInt(0).SetBytes(phoneBytes).String()
			}

		}
		user.Phone = phone
		err = userService.Get(user)
	}
	phoneInt, _ := big.NewInt(0).SetString(user.Phone, 10)
	utils.Response(c, err, userMod.User{Avatar: user.Avatar, NickName: user.NickName, Phone: hex.EncodeToString(phoneInt.Bytes()), Introduce: user.Introduce})
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
