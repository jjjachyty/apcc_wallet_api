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
	"strconv"

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

//获取兑换币种
func (AssetController) ExchangeAssets(c *gin.Context) {
	var err error
	var assets []assetMod.Asset
	mainCoin, hasMainCoin := c.GetQuery("mainCoin")
	exchangeCoin, hasExchangeCoin := c.GetQuery("exchangeCoin")
	if hasMainCoin && hasExchangeCoin {
		claims := jwt.GetClaims(c)
		if assets, err = assetService.FindExchange(claims.UUID, mainCoin, exchangeCoin); err == nil {
			fmt.Println(assets[0].Symbol, mainCoin)
			if assets[0].Symbol == mainCoin { //排序
				assets[0], assets[1] = assets[1], assets[0]
			}
		}
	}
	utils.Response(c, err, assets)
}

func (AssetController) Exchange(c *gin.Context) {
	var err error

	var amountF float64
	var assets []assetMod.Asset

	mainCoin := c.PostForm("mainCoin")
	exchangeCoin := c.PostForm("exchangeCoin")
	amount := c.PostForm("amount")
	if mainCoin != "" && exchangeCoin != "" && amount != "" {
		claims := jwt.GetClaims(c)

		// if rate, err = commonSrv.GetExchange(mainCoin, exchangeCoin); err != nil {
		if assets, err = assetService.FindExchange(claims.UUID, mainCoin, exchangeCoin); err == nil && len(assets) == 2 {
			var from = assets[0]
			var to = assets[1]

			if from.Symbol == mainCoin {
				from = assets[1]
				to = assets[0]
			}
			if amountF, err = strconv.ParseFloat(amount, 64); err == nil {
				if from.Blance > amountF {
					assetService.TransferBlance(from, to, amountF)
				} else {
					err = errors.New("余额不足")
				}
			} else {
				err = errors.New("转出金额格式错误")
			}

		}
		// }
	}
	utils.Response(c, err, assets)
}

func (AssetController) List(c *gin.Context) {
	var err error
	var assets []assetMod.Asset
	claims := jwt.GetClaims(c)
	if assets, err = assetService.Find(claims.UUID); err == nil {
		fmt.Println(assets)
		for i, asset := range assets {

			assets[i].PriceCny = commonSrv.CoinPrice[asset.NameEn]
		}
	}

	utils.Response(c, err, assets)
}

func (AssetController) Free(c *gin.Context) {
	var err error
	var free float64
	var ok bool
	coin, ok := c.GetQuery("coin")
	if ok {
		if free, ok = commonSrv.CoinFreee[coin]; !ok {
			err = errors.New("未找到币种的手续费")
		}
	} else {
		err = errors.New("缺失币种参数")
	}
	utils.Response(c, err, free)
}
