package main

import (
	"apcc_wallet_api/handler"
	"apcc_wallet_api/middlewares/cors"
	"apcc_wallet_api/router"
	"apcc_wallet_api/services/walletSrv"
	"apcc_wallet_api/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	//初始化MHC USDT 客户端
	walletSrv.InitMHCClient()
	walletSrv.InitUSDTClient()
	//开始监听消息
	handler.InitExchangeHandler()
	handler.InitTransferHandler()

	gin.DisableConsoleColor()
	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)
	wallet := gin.Default()
	wallet.Use(cors.CorsSeting())
	//lrm.Use(auth.Auth())
	//engine.Start()
	// lrm.StaticFS("/web", http.Dir("web"))
	router.WebRouter(wallet)
	//go websocket.ListenEvent()
	wallet.Run(utils.GetPort().Port)
}
