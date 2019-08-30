package exchangeCtr

import (
	"apcc_wallet_api/middlewares/jwt"
	"apcc_wallet_api/models/assetMod"
	"apcc_wallet_api/models/dimMod"
	"apcc_wallet_api/services/walletSrv"
	"apcc_wallet_api/utils"
	"encoding/json"
	"errors"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
)

type MHCExchangeController struct{}

//MHC2USDTExchange MHC 兑换USDT
func (MHCExchangeController) USDT(c *gin.Context) {
	var err = errors.New("缺少参数")
	var log assetMod.ExchangeLog
	var ethAddress common.Address
	var address string
	var tx *types.Transaction
	privateKey, hasPrivateKey := c.GetPostForm("privateKey")
	// exchangeAddress, hasExchangeAddress := c.GetPostForm("exchangeAddress")
	toAddress, hasToAddress := c.GetPostForm("toAddress")
	amountStr, hasAmount := c.GetPostForm("amount")
	log.ToCoin, log.FromCoin = "USDT", "MHC"
	log.UUID = utils.GetUUID()
	if hasPrivateKey && hasAmount && hasToAddress {

		var fromCoin, toCoin dimMod.DimCoin
		if fromCoin, err = dimCoinService.GetCoin(log.FromCoin); err == nil {
			if toCoin, err = dimCoinService.GetCoin(log.ToCoin); err == nil {
				var exchangeAddress interface{}
				if exchangeAddress, err = utils.HGet("mhc", "exchange_address"); err == nil && exchangeAddress != nil {

					if exchangeAddress, ok := exchangeAddress.(string); ok {

						exchangeRate := fromCoin.PriceCny / toCoin.PriceCny
						if bigAmount, ok := big.NewFloat(0).SetString(amountStr); ok {

							if exchangeFree, ok := dimCoinService.GetExchangeFree("MHC"); ok {
								amount, _ := bigAmount.Float64()
								usdtAmount := amount * exchangeRate
								log.FromAmount = amount
								log.User = jwt.GetClaims(c).UUID
								log.ToAddress = toAddress
								log.ToAmount = usdtAmount
								log.Free = exchangeFree
								log.FromPriceCny = fromCoin.PriceCny
								log.ToPriceCny = toCoin.PriceCny
								//处理手续费  原始金额 +手续费
								bigAmount = new(big.Float).Add(bigAmount, big.NewFloat(exchangeFree))

								amoutWeiStr := new(big.Float).Mul(bigAmount, big.NewFloat(math.Pow10(18))).Text('f', 0)
								if mhcAmount, ok := new(big.Int).SetString(amoutWeiStr, 0); ok {

									//检查账户余额是否足够
									if ethAddress, err = walletSrv.GetETHAddressByPKHex(privateKey); err == nil {

										var blance *big.Int
										if blance, err = walletSrv.GetMHCBalance(ethAddress.Hex()); err == nil {

											if blance.Cmp(new(big.Int).Add(mhcAmount, walletSrv.GetMHCGas())) > -1 {
												//足够转账
												if address, tx, err = walletSrv.SendMHCByPrivateKey(privateKey, mhcAmount, exchangeAddress); err == nil {
													log.SendTxs = tx.Hash().Hex()
													log.CreateAt = time.Now()
													log.FromAddress = address
													log.SendAddress = address
													log.Free += float64(tx.GasPrice().Int64()) * float64(tx.Gas()) / 1000000000000000000
													logJSONByts, _ := json.Marshal(log)
													//发布兑换消息
													err = utils.NsqPublish("MHC2USDT", logJSONByts) //日志

												}
											} else {
												err = errors.New("账户余额不足兑换金额(含手续费)")
											}
										}
									}
								}
							}
						}
					}
				} else {
					err = errors.New("兑换地址查询出错,请稍后再试")
				}
			}
		}
	}
	utils.Response(c, err, log)

}
