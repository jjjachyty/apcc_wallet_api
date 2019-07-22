package dimSrv

import (
	"apcc_wallet_api/models"
	"apcc_wallet_api/models/dimMod"
	"apcc_wallet_api/utils"
)

type DimCoinService struct{}

func (DimCoinService) Find() ([]dimMod.DimCoin, error) {
	coins := make([]dimMod.DimCoin, 0)
	err := models.GetBeans(&coins, dimMod.DimCoin{State: utils.STATE_ENABLE})
	return coins, err
}
