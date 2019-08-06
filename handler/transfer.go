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

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/davecgh/go-spew/spew"
)

func InitTransferHandler() {
	go utils.ReadMessage("USDT2USDT", usdt2usdt)
	go utils.ReadMessage("UpdateTransfer", updateTransfer)
	loadAddress()

}

var assetService assetSrv.AssetService
var dimCoinService dimSrv.DimCoinService
var userService userSrv.UserService

// usdt汇总地址
var toAddress = "0xC05dEb0C5e841Aa564f41f769335BC96D75Ade65"

func usdt2usdt(data []byte) (err error) {
	fmt.Println("收到USDT2USDT消息")
	var log = new(assetMod.TransferLog)
	if err = json.Unmarshal(data, &log); err == nil {
		var tx *types.Transaction
		amountStr := new(big.Float).Mul(big.NewFloat(log.Amount-log.Free), big.NewFloat(math.Pow10(6))).Text('f', 0)
		if amount, ok := big.NewInt(0).SetString(amountStr, 0); ok {
			if tx, err = walletSrv.SendUSDT(log.ToAddress, amount); err == nil {
				log.SendAddress = walletSrv.GetAuth().From.Hex()
				log.SendAt = time.Now()
				log.SendTxs = tx.Hash().Hex()
				log.State = utils.STATE_ENABLE
				transferLog, _ := json.Marshal(log)
				utils.NsqPublish("UpdateTransfer", transferLog)
			}
		} else {
			err = fmt.Errorf("USDT转出金额[%s]解析出错", amountStr)
		}

	}
	return err
}

func updateTransfer(data []byte) (err error) {
	fmt.Println("收到UpdateTransfer消息")
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
	go scanAddress()
}

func scanAddress() {
	if free, ok := dimCoinService.GetFree("USDT"); ok {
		maxCountID := userService.GetMaxCountID()
		for index := 0; index <= maxCountID; index++ {

			if address, err := walletSrv.GetEthAddress(uint32(index)); err == nil {

				if balance, err := walletSrv.GetUSDTBalance(address); err == nil {
					usdtAmount, _ := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1000000)).Float64()
					if usdtAmount > free {
						utils.SysLog.Debugf("地址%s余额%f大于手续费%f,开始转入", address, usdtAmount, free)
						var pk []byte
						if pk, err = walletSrv.GetAddressPrivateKey(uint32(index)); err == nil {
							var tx *types.Transaction
							if tx, err = walletSrv.SendUSDTByPrivateKey(common.Bytes2Hex(pk), toAddress, balance); err == nil {
								//生成转账记录
								if err := assetService.CreateLog(assetMod.TransferLog{UUID: utils.GetUUID(),
									Coin: "USDT", 
									SendTxs: tx.Hash().Hex(),
									Amount:usdtAmount,
									Free:0,
									SendAddress:
									}); err != nil {

								}
							}
						}

					}
				} else {
					utils.SysLog.Errorf("获取USDT_ETH %s余额出错 ", address)
				}

			} else {
				utils.SysLog.Errorf("从缓存获取USDT_ETH地址出错")
			}
		}
	} else {
		utils.SysLog.Errorf("获取USDT_ETH 手续费错误")
	}

}
