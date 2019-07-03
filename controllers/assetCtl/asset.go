package assetCtl

import (
	"apcc_wallet_api/middlewares/jwt"
	"apcc_wallet_api/models/assetMod"
	"apcc_wallet_api/models/userMod"
	"apcc_wallet_api/services/assetSrv"
	"apcc_wallet_api/services/commonSrv"
	"apcc_wallet_api/services/userSrv"
	"apcc_wallet_api/utils"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

type AssetController struct{}

var userService userSrv.UserService
var assetService assetSrv.AssetService

func (AssetController) Transfer(c *gin.Context) {
	var err = errors.New("参数缺失")
	var payPass = false
	address := c.PostForm("address")
	amount := c.PostForm("amount")
	payPasswd := c.PostForm("payPasswd")
	claims := jwt.GetClaims(c)
	if address != "" && amount != "" && payPasswd != "" {
		user := new(userMod.User)
		user.UUID = claims.UUID
		user.PayPasswd = utils.GetMD5(payPasswd)
		//检查支付密码
		if payPass, err = userService.CheckPayPasswd(user); payPass {
			//支付密码通过

		}
	}
	utils.Response(c, err, nil)
}

func (AssetController) List(c *gin.Context) {
	var err error
	var assets []assetMod.Asset
	claims := jwt.GetClaims(c)
	if assets, err = assetService.Find(claims.UUID); err == nil {
		fmt.Println(assets)
		for _, asset := range assets {
			fmt.Println("111111111", asset.NameEn)
			asset.PriceCny = commonSrv.CoinPrice[asset.NameEn]
		}
	}

	utils.Response(c, err, assets)
}
