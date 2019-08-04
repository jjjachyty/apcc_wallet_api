package exchange

import (
	"apcc_wallet_api/models/assetMod"
	"apcc_wallet_api/services/assetSrv"
	"apcc_wallet_api/utils"
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

var assetService assetSrv.AssetService

func init() {
	MHC2USDT()
}

func MHC2USDT() {
	go utils.ReadMessage("MHC2USDT", mhc2usdtHandler)
}

func mhc2usdtHandler(data []byte) (err error) {
	fmt.Println("收到MHC2USDT消息")
	var log = new(assetMod.AssetLog)
	if err = json.Unmarshal(data, &log); err == nil {
		spew.Dump(log)

		err = assetService.MHC2Coin(*log)
	}
	return err
}
