package handler

import (
	"apcc_wallet_api/models/assetMod"
	"apcc_wallet_api/services/assetSrv"
	"apcc_wallet_api/services/dimSrv"
	"apcc_wallet_api/services/userSrv"
	"apcc_wallet_api/services/walletSrv"
	"apcc_wallet_api/utils"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/davecgh/go-spew/spew"
)

var assetService assetSrv.AssetService
var dimCoinService dimSrv.DimCoinService
var userService userSrv.UserService
var confirmationBlockNum int64 = 3

// usdt汇总地址
// var toAddress = "0xD9B459ca8a6b03f034Fac080B6e3Ac3F830f4C9c"

func InitTransferHandler() {
	go utils.ReadMessage("USDT2USDT", usdt2usdt)
	go utils.ReadMessage("UpdateTransfer", updateTransfer)
	// go utils.ReadMessage("TransferSuccess", transferSuccess)

	go loadAddress()

}

func usdt2usdt(data []byte) (err error) {
	var log = new(assetMod.TransferLog)
	if err = json.Unmarshal(data, &log); err == nil {
		var tx *types.Transaction
		amountStr := new(big.Float).Mul(big.NewFloat(log.Amount-log.Free), big.NewFloat(math.Pow10(6))).Text('f', 0)
		if amount, ok := big.NewInt(0).SetString(amountStr, 0); ok {
			//检查热钱包USDT数量
			var bal *big.Int
			if bal, err = walletSrv.GetHotWalletUSDTBalance(); err == nil {
				if bal.Cmp(amount) > -1 {
					if tx, err = walletSrv.SendUSDT(log.ToAddress, amount); err == nil {
						log.SendAddress = walletSrv.GetAuth().From.Hex()
						log.SendAt = time.Now()
						log.SendTxs = tx.Hash().Hex()
						log.State = utils.STATE_ENABLE
						transferLog, _ := json.Marshal(log)
						utils.NsqPublish("UpdateTransfer", transferLog)
					}
				} else {
					err = fmt.Errorf("USDT提款失败,热钱包[%s]余额[%s]不足转出[%s]请尽快充值", walletSrv.GetETHHotWalletAddress(), bal.String(), amount.String())
					utils.SysLog.Errorln(err)

				}
			}
		} else {
			err = fmt.Errorf("USDT转出金额[%s]解析出错", amountStr)
		}

	}
	return err
}

func updateTransfer(data []byte) (err error) {
	var log = new(assetMod.TransferLog)
	if err = json.Unmarshal(data, &log); err == nil {
		spew.Dump(log)
		err = assetService.UpdateTransferLog(*log)
	}
	return err
}

func loadAddress() {
	var err error
	assets := make([]assetMod.Asset, 0)
	if err = assetService.Find(&assets, assetMod.Asset{Symbol: "USDT", BaseOn: "ETH"}); err == nil {
		utils.SysLog.Debugln("加载所有的USDT_ETH账户", len(assets))
		for _, asset := range assets {
			utils.SAdd("USDT-ETH", asset.Address)
		}
	}
	scanAddress()
	scanTxs()
	t1 := time.NewTimer(time.Second * 15)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 15)
			go scanAddress()
			go scanTxs()
		}
	}
}

func scanAddress() {
	if free, ok := dimCoinService.GetFree("USDT"); ok {

		maxCountID := userService.GetMaxCountID()
		// utils.SysLog.Debugf("当前一共有%d个地址", maxCountID)
		//检查热钱包余钱是否足够
		var gas = walletSrv.GetGas()
		if hotWalletUSDTBalance, err := walletSrv.GetETHHotWalletAddressBalance(); err == nil {
			var need = new(big.Int).Mul(gas, big.NewInt(maxCountID))
			if hotWalletUSDTBalance.Cmp(need) > -1 {
				for index := int64(1); index <= maxCountID; index++ {
					if address, err := walletSrv.GetEthAddress(uint32(index)); err == nil {
						if txs, err := utils.HMGet("transfer", address); err == nil && len(txs) == 1 {
							if balance, err := walletSrv.GetUSDTBalance(address); err == nil {
								// utils.SysLog.Debugf("地址%s  币种 %s 余额 %s", address, "USDT", balance.String())

								usdtAmount, _ := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1000000)).Float64()
								if usdtAmount > free {
									//检查ETH是否足够
									var ethBal *big.Int

									if ethBal, err = walletSrv.GetETHBalance(address); err == nil {
										if ethBal.Cmp(gas) > -1 { //转账足够

											utils.SysLog.Debugf("地址%s余额%f大于手续费%f,开始转入", address, usdtAmount, free)
											var pk []byte
											if pk, err = walletSrv.GetETHAddressPrivateKey(uint32(index)); err == nil {
												var tx *types.Transaction
												var auth *bind.TransactOpts
												if tx, auth, err = walletSrv.SendUSDTByPrivateKey(common.Bytes2Hex(pk), walletSrv.GetETHClodWalletAddress(), balance); err == nil {
													usdtCoin, _ := dimCoinService.GetCoin("USDT")
													transfer := assetMod.TransferLog{UUID: utils.GetUUID(),
														Coin:        "USDT",
														SendTxs:     tx.Hash().Hex(),
														Amount:      usdtAmount,
														FromAddress: auth.From.Hex(),
														PriceCny:    usdtCoin.PriceCny,
														Free:        0,
														PayType:     assetSrv.PAY_TYPE_TRANSFER_ADD_IN,
														SendAddress: walletSrv.GetAuth().From.Hex(),
														SendAt:      time.Now(),
														ToAddress:   address,
														CreateAt:    time.Now(),
													}
													transferByts, _ := json.Marshal(transfer)
													utils.HSet("transfer", address, transferByts)
													//生成转账记录
													if err := assetService.CreateLog(transfer); err != nil {
														if transferByts, err := json.Marshal(transfer); err == nil {
															utils.NsqPublish("InsertTransfer", transferByts)
														} else {
															utils.SysLog.Errorf("InsertTransfer||%s", transferByts)
														}
													}
												} else {
													utils.SysLog.Errorf("转账错误%v", err)
													//短信通知
												}
											} else {
												utils.SysLog.Errorf("获取账户地址私钥错误%v", err)
											}
										} else {
											//地址ETH不足帐,开始充值
											utils.SysLog.Warningf("地址%s转账USDT余额[%s]不足[%s]转账,开始自动充值%s", address, ethBal.String(), gas.String(), gas.String())

											if err = walletSrv.SendETH(address, gas); err == nil {
												utils.SysLog.Errorf("为转账自动充值eth失败,请手动充值%s至%s %v", gas.String(), address, err)

											}
										}
									}
								}
							} else {
								utils.SysLog.Errorf("获取USDT_ETH %s余额出错 %v", address,err)
							}
						} else {
							utils.SysLog.Warnf("地址[%s]获取出错 %v %v", address, err, txs)
						}
					} else {
						utils.SysLog.Errorf("从获取USDT_ETH地址出错[accountID=%d]", index)
					}
				}
			} else {
				utils.SysLog.Errorf("热钱包eth余额不足,手动充值[%s]至%s或更换热钱包", big.NewInt(0).Sub(need, hotWalletUSDTBalance).String(), walletSrv.GetETHHotWalletAddress().Hex())
				//短信通知热钱包充值eth
			}
		} else {
			utils.SysLog.Errorf("获取热钱包余额错误")
		}
	} else {
		utils.SysLog.Errorf("获取USDT_ETH 手续费错误")
	}

}

//扫描所有的转账
func scanTxs() {
	if allTransfer, err := utils.HGetAll("transfer"); err == nil && len(allTransfer) > 0 {
		fmt.Println("有正在转移的交易")
		for address, transferStr := range allTransfer {

			var log assetMod.TransferLog
			if err = json.Unmarshal([]byte(transferStr), &log); err == nil {
				if receipt, err := walletSrv.GetETHTransactionReceipt(log.SendTxs); err == nil && (receipt.BlockNumber.Int64()+confirmationBlockNum) < walletSrv.GetETHLastBlockNum() {

					//3次区块确认,开始转账
					utils.SysLog.Debugf("地址 %s 交易 %s 超3次接受 开始入账 %f", log.ToAddress, log.SendTxs, log.Amount)
					if err = assetService.AddCoin(log); err == nil {
						utils.HDel("transfer", address)
					} else {
						utils.SysLog.Errorf("入账%s失败 %v", log.Coin, err)
					}

				}
			}

		}

	} else {
		utils.SysLog.Errorf("无缓存的交易 %v", err)
	}
}

// func transferSuccess(data []byte) error {
// 	return assetService.UpdateTransferLog(assetMod.TransferLog{UUID: string(data), State: utils.STATE_ENABLE})
// }
