package assetCtl

import (
	"apcc_wallet_api/middlewares/jwt"
	"apcc_wallet_api/models/assetMod"
	"apcc_wallet_api/models/userMod"
	"apcc_wallet_api/services/assetSrv"
	"apcc_wallet_api/services/dimSrv"
	"apcc_wallet_api/services/userSrv"
	"apcc_wallet_api/services/walletSrv"
	"apcc_wallet_api/utils"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"

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

//Exchange BDCoin 转移到BlockCoin
func (AssetController) Exchange(c *gin.Context) {
	var err error
	var ok bool

	var log assetMod.AssetLog

	if err = c.BindJSON(&log); err == nil {
		// //检查支付密码
		// // if payPass, err = userService.CheckPayPasswd(user); payPass {
		log.UUID = utils.GetUUID()
		log.FromUser = jwt.GetClaims(c).UUID
		fromPriceCny := dimCoinService.GetCoin(log.FromCoin).PriceCny
		toPriceCny := dimCoinService.GetCoin(log.ToCoin).PriceCny
		exchangeRate := fromPriceCny / toPriceCny
		log.FromPriceCny = fromPriceCny
		log.ToPriceCny = toPriceCny
		exchangeFreeRate, _ := dimCoinService.GetExchangeFreeRate(log.FromCoin)
		log.PayType = int(assetSrv.PAY_TYPE_EXCHANGE)
		// // } else {
		// // 	err = errors.New("支付密码不正确")
		// // }
		fmt.Println("txs", log.ExchangeTxs)
		fmt.Println("log", log)
		fmt.Println(log.FromCoin + log.ToCoin)
		switch exgType := log.FromCoin + log.ToCoin; exgType {
		case "MHCUSDT":
			fmt.Println(walletSrv.GetTxsByHashHex(log.ExchangeTxs))
			log.ToAmount = (log.FromAmount * (1 - exchangeFreeRate)) * exchangeRate
			log.Free = log.FromAmount * exchangeFreeRate

			if tx, _, err := walletSrv.GetTxsByHashHex(log.ExchangeTxs); err == nil && tx != nil {
				log.State = utils.STATE_ENABLE
				err = assetService.MHC2Coin(log)
			}

		case "USDTMHC":
			var toAmount *big.Int
			mhcAmount := (log.FromAmount * (1 - exchangeFreeRate)) * exchangeRate * 1000000
			mhcAmountStr := strconv.FormatFloat(mhcAmount, 'f', 0, 64)
			toAmount, ok = big.NewInt(0).SetString(mhcAmountStr, 0)
			toAmount = new(big.Int).Mul(toAmount, big.NewInt(1000000000000))
			fmt.Println(toAmount.String())

			log.ToAmount = mhcAmount
			log.Free = (log.FromAmount * exchangeFreeRate)
			fmt.Println("ToAmount", ok, toAmount.String(), log.ToAmount)
			if ok {
				if address, txs, err := walletSrv.SendMHC(toAmount, log.ToAddress); err == nil && txs != "" {
					log.ToAddress = address
					log.SendTxs = txs
					log.State = utils.STATE_ENABLE
					err = assetService.Coin2MHC(log)
				}
			} else {
				err = errors.New("兑换金额出错")
			}

		}
	}

	utils.Response(c, err, log)
}

//Exchange BDCoin 转移到BlockCoin
func (AssetController) AssetLogUpdate(c *gin.Context) {
	var err error

	var log assetMod.AssetLog

	if err = c.BindJSON(&log); err == nil {

		err = assetService.UpdateLogs(log)

	}

	utils.Response(c, err, log)
}

//AssetsLogs 获取我的转账记录
func (AssetController) ExchangeList(c *gin.Context) {
	var err error
	var exchange = new(assetMod.Exchange)

	var page = utils.GetPageData(c)
	if err = c.ShouldBindQuery(exchange); err == nil {
		exchange.User = jwt.GetClaims(c).UUID
		err = exchangeService.GetExchanges(page, exchange)
	}
	utils.Response(c, err, page)
}

func (AssetController) Transfer(c *gin.Context) {
	var err = errors.New("参数缺失")
	var payPass bool
	var assets []assetMod.Asset
	var assetLog assetMod.AssetLog
	log, haslog := c.Get("log")
	payPasswd, hasPassword := c.Get("payPasswd")

	claims := jwt.GetClaims(c)

	if haslog && hasPassword {

		//检查支付密码
		user := new(userMod.User)
		user.UUID = claims.UUID
		user.PayPasswd = utils.GetMD5(payPasswd.(string))
		//检查支付密码
		if payPass, err = userService.CheckPayPasswd(user); payPass {
			if err = json.Unmarshal([]byte(log.(string)), &assetLog); err == nil {
				var fromPriceCny = dimCoinService.GetCoin(assetLog.FromCoin).PriceCny
				var freeRate, _ = dimCoinService.GetFreeRate(assetLog.FromCoin)
				assetLog.FromPriceCny = fromPriceCny
				assetLog.Free = assetLog.FromAmount * freeRate
				//支付密码通过
				switch assetLog.PayType {
				case assetSrv.PAY_TYPE_TRANSFER_INNER: //内部转账

					// if assets, err = assetService.FindInnerTransfer(assetMod.Asset{UUID: claims.UUID, Address: fromAddress, Symbol: symbol}, assetMod.Asset{Address: toAddress}); err == nil && len(assets) == 2 {
					// if assets[0].Blance >= amountF {
					err = assetService.Send(assetLog)
					// } else {
					// 	err = errors.New("余额不足")
					// }

					// } else {
					// 	err = errors.New("平台未找到该地址,请核实")
					// }
				case assetSrv.PAY_TYPE_TRANSFER_OUTER: // 转出平台

					if freeRate, ok := dimCoinService.GetFreeRate(assetLog.FromCoin); ok {
						assetLog.ToAmount = assetLog.FromAmount * (1 - freeRate)
						//转账到外部地址
						err = assetService.Send(assetLog)
					} else {
						err = fmt.Errorf("未找到%s的手续费", assetLog.FromCoin)
					}
				}

			} else {
				err = errors.New("支付密码不正确")
			}
		}
	}
	utils.Response(c, err, assets)
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
	var assetsLog = new(assetMod.AssetLog)

	var page = utils.GetPageData(c)
	if err = c.ShouldBindQuery(assetsLog); err == nil {
		assetsLog.FromUser = jwt.GetClaims(c).UUID
		err = assetService.GetLogs(page, assetsLog)
	}
	utils.Response(c, err, page)
}

func (AssetController) TransferFreeRate(c *gin.Context) {
	var err error
	var free float64
	symbol, ok := c.GetQuery("symbol")
	if ok {
		if free, ok = dimCoinService.GetFreeRate(symbol); !ok {
			err = errors.New("未找到币种的手续费")
		}
	} else {
		err = errors.New("缺失币种参数")
	}
	utils.Response(c, err, free)
}

//获取兑换费用
func (AssetController) ExchangeFreeRate(c *gin.Context) {
	var err error
	var exchangeFree float64
	mainCoin, ok := c.GetQuery("mainCoin")
	// exchange, exchangeOk := c.GetQuery("exchangeCoin")

	if ok {
		if exchangeFree, ok = dimCoinService.GetExchangeFreeRate(mainCoin); !ok {
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
