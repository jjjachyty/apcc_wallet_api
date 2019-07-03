package assetSrv

import (
	"apcc_wallet_api/models"
	"apcc_wallet_api/models/assetMod"
	"apcc_wallet_api/services/userSrv"
)

const (
	HMC_COIN_SYMBOL  = "MHC"
	USDT_COIN_SYMBOL = "USDT"
)

var (
	assetsSQL = "select a.*,b.name_en ,b.name_cn FROM asset a LEFT JOIN dim_coin b on  a.symbol = b.symbol where a.uuid=?"
)

type AssetService struct{}

var userService userSrv.UserService

func TransferCoin(symbol string, fromAddress string, amount float64, toAddress string) {
	switch symbol {
	case HMC_COIN_SYMBOL:

	case USDT_COIN_SYMBOL:
	}
}

func (AssetService) Create(assets []assetMod.Asset) error {

	return models.Create(&assets)
}

func (AssetService) Find(uuid string) ([]assetMod.Asset, error) {
	var assets = make([]assetMod.Asset, 0)
	return assets, models.SQLBeans(&assets, assetsSQL, uuid)
}
