package transferCtr

import (
	"apcc_wallet_api/middlewares/jwt"
	"apcc_wallet_api/models/assetMod"
	"apcc_wallet_api/models/userMod"
	"apcc_wallet_api/services/assetSrv"
	"apcc_wallet_api/services/dimSrv"
	"apcc_wallet_api/services/userSrv"
	"apcc_wallet_api/utils"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type USDTTransferController struct{}

var userService userSrv.UserService
var assetService assetSrv.AssetService
var exchangeService assetSrv.ExchangeService
var dimCoinService dimSrv.DimCoinService

func (USDTTransferController) Transfer(c *gin.Context) {

	var err error

	var transferLog assetMod.TransferLog
	fromAddress, hasFromAddress := c.GetPostForm("fromAddress")
	toAddress, hasToAddress := c.GetPostForm("toAddress")
	amountStr, hasAmount := c.GetPostForm("amount")
	payPasswd, hasPassword := c.GetPostForm("password")

	claims := jwt.GetClaims(c)

	if hasToAddress && hasAmount && hasPassword && hasFromAddress {
		if fromAddress != toAddress {

			//检查支付密码
			user := new(userMod.User)
			user.UUID = claims.UUID
			user.PayPasswd = utils.GetMD5(payPasswd)
			//检查支付密码
			var payPass bool
			if payPass, err = userService.CheckPayPasswd(user); err == nil && payPass {
				transferLog.UUID = utils.GetUUID()
				transferLog.FromUser = claims.UUID

				transferLog.ToAddress, transferLog.FromAddress = toAddress, fromAddress
				transferLog.Coin = "USDT"
				transferLog.PriceCny = dimCoinService.GetCoin(transferLog.Coin).PriceCny
				if free, ok := dimCoinService.GetFree(transferLog.Coin); ok {
					transferLog.Free = free

					var amount float64
					if amount, err = strconv.ParseFloat(amountStr, 0); err == nil {
						transferLog.Amount = amount + free

						//检查地址是否是内部地址
						var asset = new(assetMod.Asset)
						asset.Address = toAddress
						if err = assetService.GetBean(asset); err == nil {
							if asset.UUID != "" { //内部转账
								transferLog.PayType = assetSrv.PAY_TYPE_TRANSFER_INNER
								transferLog.Free = 0 //内部转账
								transferLog.State = utils.STATE_ENABLE
								err = assetService.SendInner(transferLog)
							} else { //外部转账
								transferLog.PayType = assetSrv.PAY_TYPE_TRANSFER_OUTER
								var assetLogJSONByts []byte
								if assetLogJSONByts, err = json.Marshal(transferLog); err == nil {
									if err = assetService.SendOuter(transferLog); err == nil {
										utils.SysLog.Debugln("开始发送USDT2USDT消息")
										utils.NsqPublish("USDT2USDT", assetLogJSONByts)
									}
								}
							}
						}
					} else {
						err = errors.New("金额格式错误")
					}
				} else {
					err = fmt.Errorf("获取%s手续费失败", transferLog.Coin)
				}
			} else {
				err = errors.New("支付密码不正确")
			}

		} else {
			err = errors.New("转入转出地址不能相同")
		}
	}
	utils.Response(c, err, transferLog)
}
