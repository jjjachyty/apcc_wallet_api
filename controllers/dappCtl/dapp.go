package dappCtl

import (
	"apcc_wallet_api/models/dappMod"
	"apcc_wallet_api/services/dappSrv"
	"apcc_wallet_api/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

//DappController Dapp控制器
type DappController struct{}

var dappService dappSrv.DappService

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
