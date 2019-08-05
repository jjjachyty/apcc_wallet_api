package handler

import (
	"apcc_wallet_api/models/assetMod"
	"apcc_wallet_api/services/walletSrv"
	"apcc_wallet_api/utils"
	"encoding/json"
	"fmt"
)

func init() {
	go utils.ReadMessage("USDT2USDT", USDT2USDT)

}

func USDT2USDT(data []byte) (err error) {
	fmt.Println("收到USDT2USDT消息")
	var log = new(assetMod.TransferLog)
	if err = json.Unmarshal(data, &log); err == nil {
		walletSrv.SendUSDT()
	}
}
