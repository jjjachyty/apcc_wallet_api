package exchangeCtr

import (
	"apcc_wallet_api/middlewares/jwt"
	"apcc_wallet_api/models/assetMod"
	"apcc_wallet_api/models/dimMod"
	"apcc_wallet_api/models/userMod"
	"apcc_wallet_api/services/assetSrv"
	"apcc_wallet_api/services/dimSrv"
	"apcc_wallet_api/services/userSrv"
	"apcc_wallet_api/utils"
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"github.com/gin-gonic/gin"
)

type USDTExchangeController struct{}

var dimCoinService dimSrv.DimCoinService
var userService userSrv.UserService
var exchangeService assetSrv.ExchangeService

//MHC2USDTExchange MHC 兑换USDT
func (USDTExchangeController) MHC(c *gin.Context) {
	var err error
	var log assetMod.ExchangeLog

	fromAddress, hasFromAddress := c.GetPostForm("fromAddress")
	toAddress, hasToAddress := c.GetPostForm("toAddress")
	amountStr, hasAmount := c.GetPostForm("amount")
	passwordStr, hasPassword := c.GetPostForm("password")
	log.ToCoin, log.FromCoin = "MHC", "USDT"
	log.UUID = utils.GetUUID()
	if hasFromAddress && hasAmount && hasToAddress && hasPassword {
		claims := jwt.GetClaims(c)
		var ok bool
		if ok, err = userService.CheckPayPasswd(&userMod.User{UUID: claims.UUID, PayPasswd: utils.GetMD5(passwordStr)}); err == nil && ok {
			var fromCoin, toCoin dimMod.DimCoin
			if fromCoin, err = dimCoinService.GetCoin(log.FromCoin); err == nil {
				if toCoin, err = dimCoinService.GetCoin(log.ToCoin); err == nil {

					exchangeRate := fromCoin.PriceCny / toCoin.PriceCny

					if bigAmount, ok := big.NewFloat(0).SetString(amountStr); ok {
						if exchangeFree, ok := dimCoinService.GetExchangeFree("USDT"); ok {
							amount, _ := bigAmount.Float64()
							mhcAmount := (amount * exchangeRate) - exchangeFree
							log.FromAmount = amount
							log.User = claims.UUID
							log.FromAddress = fromAddress
							log.ToAddress = toAddress
							log.ToAmount = mhcAmount
							log.Free = exchangeFree
							log.FromPriceCny = fromCoin.PriceCny
							log.ToPriceCny = toCoin.PriceCny
							log.CreateAt = time.Now()
							var pubData []byte
							if pubData, err = json.Marshal(log); err == nil {
								if err = exchangeService.SubCoin(log); err == nil {
									//发布兑换消息
									utils.NsqPublish("USDT2MHC", pubData)
								}

							}

						}
					}
				}
			}
		} else {
			err = errors.New("密码校验失败")
		}
	}
	utils.Response(c, err, log)

}
