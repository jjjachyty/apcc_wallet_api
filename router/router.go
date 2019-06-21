package router

import (
	"apcc_wallet/controllers/common"

	"github.com/gin-gonic/gin"
)

func WebRouter(router *gin.Engine) {

	// RestFUl 路由
	v1 := router.Group("/api/wallet/v1/") //版本
	{

		com := v1.Group("/com") //参数模块
		{
			com.Any("/sms", common.SMSController{}.Controller)
			com.Any("/captcha", common.CaptchaController{}.Controller)
		}
	}
}
