package dimCtr

import (
	"apcc_wallet_api/models/dimMod"
	"apcc_wallet_api/services/dimSrv"
	"apcc_wallet_api/utils"

	"github.com/gin-gonic/gin"
)

type DimConfigController struct{}

func (DimConfigController) AddOrUpdate(c *gin.Context) {
	var config dimMod.DimConfig
	var err error
	if err = c.Bind(&config); err == nil {
		dimSrv.AddOrUpdateConfig(config)
	}
	utils.Response(c, err, config)
}
