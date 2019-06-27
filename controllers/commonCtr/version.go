package commonCtr

import (
	"apcc_wallet_api/services/commonSrv"
	"apcc_wallet_api/utils"

	"github.com/gin-gonic/gin"
)

func GetMaxVersion(c *gin.Context) {
	version, err := commonSrv.GetLastVersion()
	utils.Response(c, err, version)
}
