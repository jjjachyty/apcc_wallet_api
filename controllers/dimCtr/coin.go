package dimCtr

import (
	"apcc_wallet_api/services/dimSrv"

	"apcc_wallet_api/utils"

	"github.com/gin-gonic/gin"
)

//DappController Dapp控制器
type DimCoinController struct{}

var dimCoinService dimSrv.DimCoinService

//All 查询COin码表
func (DimCoinController) All(c *gin.Context) {
	result, err := dimCoinService.Find()

	utils.Response(c, err, result)
}
