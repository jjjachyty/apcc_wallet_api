package assetCtr

import (
	"apcc_wallet_api/middlewares/jwt"
	"apcc_wallet_api/models/assetMod"
	"apcc_wallet_api/services/assetSrv"
	"apcc_wallet_api/services/dimSrv"
	"apcc_wallet_api/services/userSrv"
	"apcc_wallet_api/utils"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

type AssetController struct{}

var userService userSrv.UserService
var assetService assetSrv.AssetService
var exchangeService assetSrv.ExchangeService
var dimCoinService dimSrv.DimCoinService

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
			if assets[0].Symbol == exchangeCoin { //排序
				assets[0], assets[1] = assets[1], assets[0]
			}
		}
	}
	utils.Response(c, err, assets)
}

// //Exchange BDCoin 转移到BlockCoin
// func (AssetController) Exchange(c *gin.Context) {
// 	var err error
// 	var ok bool

// 	var log assetMod.AssetLog

// 	if err = c.BindJSON(&log); err == nil {
// 		// //检查支付密码
// 		// // if payPass, err = userService.CheckPayPasswd(user); payPass {
// 		log.UUID = utils.GetUUID()
// 		log.FromUser = jwt.GetClaims(c).UUID
// 		fromPriceCny := dimCoinService.GetCoin(log.FromCoin).PriceCny
// 		toPriceCny := dimCoinService.GetCoin(log.ToCoin).PriceCny
// 		exchangeRate := fromPriceCny / toPriceCny
// 		log.FromPriceCny = fromPriceCny
// 		log.ToPriceCny = toPriceCny
// 		exchangeFree, _ := dimCoinService.GetExchangeFree(log.FromCoin)
// 		log.PayType = int(assetSrv.PAY_TYPE_EXCHANGE)
// 		// // } else {
// 		// // 	err = errors.New("支付密码不正确")
// 		// // }

// 		switch exgType := log.FromCoin + log.ToCoin; exgType {
// 		case "MHCUSDT":

// 			log.ToAmount = (log.FromAmount - exchangeFree) * exchangeRate
// 			log.Free = exchangeFree

// 			if tx, _, err := walletSrv.GetTxsByHashHex(log.ExchangeTxs); err == nil && tx != nil {
// 				log.State = utils.STATE_ENABLE
// 				err = assetService.MHC2Coin(log)
// 			}

// 		case "USDTMHC":
// 			var toAmount *big.Int
// 			mhcAmount := (log.FromAmount - exchangeFree) * exchangeRate * 1000000
// 			mhcAmountStr := strconv.FormatFloat(mhcAmount, 'f', 0, 64)
// 			toAmount, ok = big.NewInt(0).SetString(mhcAmountStr, 0)
// 			toAmount = new(big.Int).Mul(toAmount, big.NewInt(1000000000000))
// 			fmt.Println(toAmount.String())

// 			log.ToAmount = mhcAmount
// 			log.Free = exchangeFree
// 			fmt.Println("ToAmount", ok, toAmount.String(), log.ToAmount)
// 			if ok {
// 				if address, txs, err := walletSrv.SendMHC(toAmount, log.ToAddress); err == nil && txs != "" {
// 					log.ToAddress = address
// 					log.SendTxs = txs
// 					log.State = utils.STATE_ENABLE
// 					err = assetService.Coin2MHC(log)
// 				}
// 			} else {
// 				err = errors.New("兑换金额出错")
// 			}

// 		}
// 	}

// 	utils.Response(c, err, log)
// }

//AssetsLogs 获取我的转账记录
func (AssetController) ExchangeList(c *gin.Context) {
	var err error
	var exchange = new(assetMod.ExchangeLog)

	var page = utils.GetPageData(c)
	if err = c.ShouldBindQuery(exchange); err == nil {
		exchange.User = jwt.GetClaims(c).UUID
		err = exchangeService.GetExchanges(page, exchange)
	}
	utils.Response(c, err, page)
}

func (AssetController) List(c *gin.Context) {
	var err error
	var assets []assetMod.Asset
	var cond assetMod.Asset

	if err = c.ShouldBindQuery(&cond); err == nil {
		cond.UUID = jwt.GetClaims(c).UUID
		err = assetService.Find(&assets, cond) //; err == nil {
		// 	for i, asset := range assets {

		// 		assets[i].PriceCny = dimCoinService.GetCoin(asset.Symbol).PriceCny
		// 	}
		// }
	}
	utils.Response(c, err, assets)
}

//AssetsLogs 获取我的转账记录
func (AssetController) Orders(c *gin.Context) {
	var err error
	var transferLog = new(assetMod.TransferLog)

	var page = utils.GetPageData(c)
	if err = c.ShouldBindQuery(transferLog); err == nil {
		transferLog.FromUser = jwt.GetClaims(c).UUID
		err = assetService.GetLogs(page, transferLog)
	}
	utils.Response(c, err, page)
}

func (AssetController) TransferFree(c *gin.Context) {
	var err error
	var free float64
	symbol, ok := c.GetQuery("symbol")
	if ok {
		if free, ok = dimCoinService.GetFree(symbol); !ok {
			err = errors.New("未找到币种的手续费")
		}
	} else {
		err = errors.New("缺失币种参数")
	}
	utils.Response(c, err, free)
}

//获取兑换费用
func (AssetController) ExchangeFree(c *gin.Context) {
	var err error
	var exchangeFree float64
	mainCoin, ok := c.GetQuery("mainCoin")
	// exchange, exchangeOk := c.GetQuery("exchangeCoin")

	if ok {
		if exchangeFree, ok = dimCoinService.GetExchangeFree(mainCoin); !ok {
			err = errors.New("币种不存在")
		}
	} else {
		err = errors.New("缺失币种参数")
	}
	utils.Response(c, err, exchangeFree)
}

//获取兑换费用
func (AssetController) GetExchangeRate(c *gin.Context) {
	var err error
	var exchangeFree float64
	mainCoin, hasMainCoin := c.GetQuery("mainCoin")
	exchangeCoin, hasExchangeCoin := c.GetQuery("exchangeCoin")
	if hasMainCoin && hasExchangeCoin {
		if exchangeFree, err = dimCoinService.GetExchangeRate(mainCoin, exchangeCoin); err == nil {

		}
	} else {
		err = errors.New("缺失币种参数")
	}
	utils.Response(c, err, exchangeFree)
}
