package dappCtr

import (
	"apcc_wallet_api/middlewares/jwt"
	"apcc_wallet_api/models/dappMod"
	"apcc_wallet_api/services/dappSrv"
	"apcc_wallet_api/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

//DappController Dapp控制器
type DappController struct{}

var dappService dappSrv.DappService
var dappUseLogService dappSrv.DappUseLogService

//Page 分页查询 dapp
func (DappController) Page(c *gin.Context) {
	var err error
	var dapp = make([]dappMod.Dapp, 0)

	var page = utils.GetPageData(c)
	page.Rows = &dapp
	fmt.Println("value", c.Request.URL.Query())
	err = dappService.Page(page, c.Request.URL.Query())

	utils.Response(c, err, page)
}

//Page 分页查询 dapp
func (DappController) Main(c *gin.Context) {
	var err error
	var dapp = make([]dappMod.Dapp, 0)

	user, ok := c.GetQuery("user")
	if ok && user != "" { //已登录的用户返回最近使用的dapp
		if dapp, err = dappUseLogService.RecentUsed(user); len(dapp) == 0 {
			dapp, err = dappService.AdvancedFind("score", 6, "state=?", utils.STATE_ENABLE)
		}
	} else { //未登录的用户返回推荐Dapp
		dapp, err = dappService.AdvancedFind("score", 6, "state=?", utils.STATE_ENABLE)
	}

	utils.Response(c, err, dapp)
}

//Page 分页查询 dapp
func (DappController) Used(c *gin.Context) {
	var err error
	var claims *jwt.CustomClaims

	token := c.Request.Header.Get("authorization")
	dapp, ok := c.GetQuery("dapp")

	if token != "" && ok { //已登录的用户返回最近使用的dapp
		claims, err = jwt.NewJWT().ParseToken(token)
		err = dappUseLogService.Add(dappMod.DappUseLog{UUID: utils.GetUUID(), User: claims.UUID, Dapp: dapp})
	}

	utils.Response(c, err, nil)
}
