package router

import (
	"apcc_wallet_api/controllers/authController"
	"apcc_wallet_api/controllers/commonController"

	"github.com/gin-gonic/gin"
)

func WebRouter(router *gin.Engine) {

	// RestFUl 路由
	v1 := router.Group("/api/wallet/v1/") //版本
	{

		com := v1.Group("/com") //参数模块
		{
			com.Any("/sms", commonController.SMSController{}.Controller)
			com.Any("/captcha", commonController.CaptchaController{}.Controller)
		}
		auth := v1.Group("/auth") //参数模块
		{
			auth.POST("/register", authController.RegisterController{}.Register)
			auth.POST("/loginwithpw", authController.RegisterController{}.LoginWithPW)
			auth.POST("/loginwithsms", authController.RegisterController{}.LoginWithSMS)
		}
	}
}
