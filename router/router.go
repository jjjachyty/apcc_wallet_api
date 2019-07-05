package router

import (
	"apcc_wallet_api/controllers/assetCtl"
	"apcc_wallet_api/controllers/commonCtr"
	"apcc_wallet_api/controllers/userCtr"
	"apcc_wallet_api/middlewares/jwt"

	"github.com/gin-gonic/gin"
)

func WebRouter(router *gin.Engine) {

	// RestFUl 路由
	v1 := router.Group("/api/wallet/v1/") //版本
	{

		com := v1.Group("/com") //参数模块
		{
			com.Any("/sms", commonCtr.SMSController{}.Controller)
			com.Any("/captcha", commonCtr.CaptchaController{}.Controller)
			com.Any("/version", commonCtr.GetMaxVersion)
		}

		auth := v1.Group("/auth") //参数模块
		{
			auth.POST("/register", userCtr.UserController{}.Register)
			auth.POST("/loginwithpw", userCtr.UserController{}.LoginWithPW)
			auth.POST("/loginwithsms", userCtr.UserController{}.LoginWithSMS)
			auth.POST("/refreshtoken", jwt.RefreshToken)

		}
		v1.Use(jwt.JWTAuth())
		user := v1.Group("/user") //参数模块
		{
			user.POST("/paypasswd", userCtr.UserController{}.PayPassword)
			user.POST("/loginpasswd", userCtr.UserController{}.LoginPassword)
			user.POST("/profile", userCtr.UserController{}.Profile)
			user.POST("/idcard", userCtr.UserController{}.IDCard)

		}
		assets := v1.Group("/assets") //参数模块
		{

			assets.GET("/all", assetCtl.AssetController{}.List)
			assets.GET("/exchangeassets", assetCtl.AssetController{}.ExchangeAssets)
			assets.POST("/exchange", assetCtl.AssetController{}.Exchange)
			assets.GET("/free", assetCtl.AssetController{}.Free)
			assets.POST("/transfer", assetCtl.AssetController{}.Transfer)
		}

		com.POST("/idcardrecognition", commonCtr.IDCardRecognition)
		test := v1.Group("/test")
		{
			test.Use(jwt.JWTAuth())
			test.POST("/jwt", jwt.GetDataByTime)
		}

	}
}
