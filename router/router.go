package router

import (
	"apcc_wallet_api/controllers/assetCtr"
	"apcc_wallet_api/controllers/commonCtr"
	"apcc_wallet_api/controllers/dappCtr"
	"apcc_wallet_api/controllers/dimCtr"
	"apcc_wallet_api/controllers/exchangeCtr"
	"apcc_wallet_api/controllers/transferCtr"
	"apcc_wallet_api/controllers/userCtr"
	"apcc_wallet_api/middlewares/jwt"

	"github.com/gin-gonic/gin"
)

func WebRouter(router *gin.Engine) {

	// RestFUl 路由
	v1 := router.Group("/api/wallet/v1/") //版本
	{

		com := v1.Group("/com") //公共
		{
			com.Any("/sms", commonCtr.SMSController{}.Controller)
			com.Any("/captcha", commonCtr.CaptchaController{}.Controller)
			com.Any("/version", commonCtr.GetMaxVersion)
			com.GET("/news", commonCtr.NewsController{}.NewsList)
			com.DELETE("/news", commonCtr.NewsController{}.RemoveNews)
			com.GET("/newsdetail", commonCtr.NewsController{}.NewsDetail)
		}

		dim := v1.Group("/dim") //纬度
		{
			dim.GET("/coins", dimCtr.DimCoinController{}.All)
			dim.POST("/config", dimCtr.DimConfigController{}.AddOrUpdate)
		}
		auth := v1.Group("/auth") //权限
		{
			auth.POST("/register", userCtr.UserController{}.Register)
			auth.GET("/checkphone", userCtr.UserController{}.CheckPhone)

			auth.POST("/loginwithpw", userCtr.UserController{}.LoginWithPW)
			auth.POST("/loginwithsms", userCtr.UserController{}.LoginWithSMS)
			auth.POST("/refreshtoken", jwt.RefreshToken)

		}
		dapp := v1.Group("/dapp") //dapp
		{

			dapp.GET("/all", dappCtr.DappController{}.Page)
			dapp.GET("/main", dappCtr.DappController{}.Main)
			dapp.GET("/used", dappCtr.DappController{}.Used)

		}

		v1.Use(jwt.JWTAuth())
		user := v1.Group("/user") //用户
		{
			user.POST("/paypasswd", userCtr.UserController{}.PayPassword)
			user.POST("/loginpasswd", userCtr.UserController{}.LoginPassword)
			user.POST("/profile", userCtr.UserController{}.Profile)
			user.POST("/idcard", userCtr.UserController{}.IDCard)

		}

		news := v1.Group("/news")
		{
			news.POST("/addorupdate", commonCtr.NewsController{}.AddOrUpdateNews)

		}
		//货币兑换
		exchange := v1.Group("/exchange")
		{
			exchange.POST("/mhc2usdt", exchangeCtr.MHCExchangeController{}.USDT)
			exchange.POST("/usdt2mhc", exchangeCtr.USDTExchangeController{}.MHC)

		}
		transfer := v1.Group("/transfer")
		{
			transfer.POST("/usdt", transferCtr.USDTTransferController{}.Transfer)

		}

		assets := v1.Group("/assets") //资产
		{

			assets.GET("/all", assetCtr.AssetController{}.List)
			assets.GET("/exchangeassets", assetCtr.AssetController{}.ExchangeAssets)
			// assets.POST("/exchange", assetCtl.AssetController{}.Exchange)
			// assets.PUT("/log", assetCtr.AssetController{}.AssetLogUpdate)
			assets.GET("/exchanges", assetCtr.AssetController{}.ExchangeList)
			assets.GET("/transferfree", assetCtr.AssetController{}.TransferFree)
			assets.GET("/exchangefree", assetCtr.AssetController{}.ExchangeFree)
			assets.GET("/exchangerate", assetCtr.AssetController{}.GetExchangeRate)

			// assets.POST("/transfer", assetCtr.AssetController{}.Transfer)
			assets.GET("/logs", assetCtr.AssetController{}.Orders)
			assets.GET("/mhclogs", assetCtr.AssetController{}.MHCOrders)

		}
		// wallet := v1.Group("/wallet") //钱包
		// {

		// 	 wallet.GET("/address", assetCtl.WalletController{}.GetAddress)
		// }
		com.POST("/idcardrecognition", commonCtr.IDCardRecognition)
		test := v1.Group("/test")
		{
			test.Use(jwt.JWTAuth())
			test.POST("/jwt", jwt.GetDataByTime)
		}

	}
}
