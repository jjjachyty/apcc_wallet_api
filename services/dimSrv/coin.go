package dimSrv

import (
	"apcc_wallet_api/models"
	"apcc_wallet_api/models/dimMod"
	"apcc_wallet_api/utils"
	"fmt"
)

type DimCoinService struct{}

var coins = make(map[string]dimMod.DimCoin)

func (DimCoinService) GetCoin(coinSymbol string) dimMod.DimCoin {
	return coins[coinSymbol]
}

func init() {
	coinRows := make([]dimMod.DimCoin, 0)
	if err := models.GetBeans(&coinRows, dimMod.DimCoin{State: utils.STATE_ENABLE}); err == nil {
		for _, coin := range coinRows {
			coins[coin.Symbol] = coin
		}
	}

}
func (DimCoinService) Find() ([]dimMod.DimCoin, error) {
	coins := make([]dimMod.DimCoin, 0)
	err := models.GetBeans(&coins, dimMod.DimCoin{State: utils.STATE_ENABLE})
	return coins, err
}

func (DimCoinService) GetExchangeRate(mainCoin, exchangeCoin string) (float64, error) {
	var exchangeRate = coins[exchangeCoin].PriceCny
	var mainRate = coins[mainCoin].PriceCny
	if exchangeRate > 0 && mainRate > 0 {
		return mainRate / exchangeRate, nil
	}
	return 0, fmt.Errorf("未找到[%s->%s]的汇率", mainCoin, exchangeCoin)
}

func (DimCoinService) GetExchangeFree(coin string) (float64, bool) {
	if coin, ok := coins[coin]; ok {
		return coin.ExchangeFree, ok
	}
	return 0, false
}
func (DimCoinService) GetFree(coin string) (float64, bool) {
	if coin, ok := coins[coin]; ok {
		return coin.TransferFree, ok
	}
	return 0, false
}
