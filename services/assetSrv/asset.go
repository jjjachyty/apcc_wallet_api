package assetSrv

import (
	"apcc_wallet_api/services/userSrv"
)

const (
	HMC_COIN_SYMBOL  = "MHC"
	USDT_COIN_SYMBOL = "USDT"
)

type AssetService struct{}

var userService userSrv.UserService

func TransferCoin(symbol string, fromAddress string, amount float64, toAddress string) {
	switch symbol {
	case HMC_COIN_SYMBOL:
		
	case USDT_COIN_SYMBOL:
	}
}
