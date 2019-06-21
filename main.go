package main

import (
	"apcc_wallet/middlewares/cors"
	"apcc_wallet/router"
	"apcc_wallet/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	//produrce.InitProdurce()
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
