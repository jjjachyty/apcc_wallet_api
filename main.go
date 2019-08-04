package main

import (
	"apcc_wallet_api/middlewares/cors"
	"apcc_wallet_api/router"
	"apcc_wallet_api/utils"

	_ "apcc_wallet_api/exchange"

	"github.com/gin-gonic/gin"
)

func main() {

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
