package dimSrv

import (
	"apcc_wallet_api/models"
	"apcc_wallet_api/models/dimMod"
	"apcc_wallet_api/utils"
	"fmt"
)

type DimCoinService struct{}

var Coins = make(map[string]dimMod.DimCoin)

func init() {
	coins := make([]dimMod.DimCoin, 0)
	if err := models.GetBeans(&coins, dimMod.DimCoin{State: utils.STATE_ENABLE}); err == nil {
		for _, coin := range coins {
			Coins[coin.Symbol] = coin
		}
	}

}
func (DimCoinService) Find() ([]dimMod.DimCoin, error) {
	coins := make([]dimMod.DimCoin, 0)
	err := models.GetBeans(&coins, dimMod.DimCoin{State: utils.STATE_ENABLE})
	return coins, err
}

func (DimCoinService) GetExchange(mainCoin, exchangeCoin string) (float64, error) {
	var exchangeRate = Coins[exchangeCoin].PriceCny
	var mainRate = Coins[mainCoin].PriceCny
	if exchangeRate > 0 && mainRate > 0 {
		return exchangeRate / mainRate, nil
	}
	return 0, fmt.Errorf("未找到[%s->%s]的汇率", mainCoin, exchangeCoin)
}

func (DimCoinService) GetExchangeUSDTFree(coin string) float64 {
	return Coins[coin].ExchangeUsdtFree
}
