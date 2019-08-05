package handler

import (
	"apcc_wallet_api/models/assetMod"
	"apcc_wallet_api/services/assetSrv"
	"apcc_wallet_api/services/walletSrv"
	"apcc_wallet_api/utils"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"time"

	"github.com/davecgh/go-spew/spew"
)

var exchangeService assetSrv.ExchangeService

func init() {
	go utils.ReadMessage("MHC2USDT", mhc2usdtHandler)
	go utils.ReadMessage("USDT2MHC", usdt2mhcHandler)

	go utils.ReadMessage("UpdateExchange", updateExchangeHandler)

}

func mhc2usdtHandler(data []byte) (err error) {
	fmt.Println("收到MHC2USDT消息")
	var log = new(assetMod.ExchangeLog)
	if err = json.Unmarshal(data, &log); err == nil {
		spew.Dump(log)
		log.State = utils.STATE_ENABLE
		log.SendAt = time.Now()
		if err = exchangeService.MHC2Coin(*log); err == nil {
			logData, _ := json.Marshal(log)
			utils.AppLog.Debugf("HMC2USDT||%s", logData)
			utils.NsqPublish("UpdateExchange", logData)
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
				fmt.Println("usdt2mhcHandler", address, txs, err)
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
