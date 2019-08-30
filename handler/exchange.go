package handler

import (
	"apcc_wallet_api/models/assetMod"
	"apcc_wallet_api/services/assetSrv"
	"apcc_wallet_api/services/walletSrv"
	"apcc_wallet_api/utils"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/davecgh/go-spew/spew"
)

var exchangeService assetSrv.ExchangeService

func InitExchangeHandler() {
	go utils.ReadMessage("MHC2USDT", mhc2usdtHandler)
	go utils.ReadMessage("USDT2MHC", usdt2mhcHandler)

	go utils.ReadMessage("UpdateExchange", updateExchangeHandler)

}

func mhc2usdtHandler(data []byte) (err error) {
	fmt.Println("收到MHC2USDT消息")
	var log = new(assetMod.ExchangeLog)
	var recpt *types.Receipt
	if err = json.Unmarshal(data, &log); err == nil {

		if recpt, err = walletSrv.GetMHCTransactionReceipt(log.SendTxs); err == nil {
			if recpt.BlockNumber.Int64()+3 < walletSrv.GetMHCLastBlockNum() { //超过3个区块接受则转账
				log.State = utils.STATE_ENABLE
				log.SendAt = time.Now()
				if err = exchangeService.AddCoin(*log); err == nil {
					logData, _ := json.Marshal(log)
					utils.AppLog.Debugf("HMC2USDT||%s", logData)
					utils.NsqPublish("UpdateExchange", logData)
				}
			} else {
				return errors.New("等待更多的区块接受")
			}
		}

	}
	return err
}

func usdt2mhcHandler(data []byte) (err error) {
	fmt.Println("收到USDT2MHC消息")
	var log = new(assetMod.ExchangeLog)
	if err = json.Unmarshal(data, &log); err == nil {
		spew.Dump(log)
		mhcAmountStr := new(big.Float).Mul(big.NewFloat(log.ToAmount), big.NewFloat(math.Pow10(18))).Text('f', 0)
		if mhcAmount, ok := big.NewInt(0).SetString(mhcAmountStr, 0); ok {

			var address, txs string
			if address, txs, err = walletSrv.SendMHC(mhcAmount, log.ToAddress); err == nil {
				log.SendAddress = address
				log.SendTxs = txs
				log.State = utils.STATE_ENABLE

				log.SendAt = time.Now()
				logData, _ := json.Marshal(log)
				utils.AppLog.Debugf("HMC2USDT||%s", logData)

				utils.NsqPublish("UpdateExchange", logData)

			}
		}

	}
	return err
}

func updateExchangeHandler(data []byte) (err error) {
	fmt.Println("收到UpdateExchange消息")
	var log = new(assetMod.ExchangeLog)
	if err = json.Unmarshal(data, &log); err == nil {
		spew.Dump(log)
		err = exchangeService.Update(log)
	}
	return err
}
