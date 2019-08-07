package dimSrv

import (
	"apcc_wallet_api/models"
	"apcc_wallet_api/models/dimMod"
	"apcc_wallet_api/utils"
	"encoding/json"
	"fmt"
)

type DimCoinService struct{}

var coins = make(map[string]dimMod.DimCoin)

func (DimCoinService) GetCoin(coinSymbol string) (dimMod.DimCoin, error) {
	var coin dimMod.DimCoin
	var coinBytes string
	var err error
	if coinBytes, err = utils.Get("dimcoin:" + coinSymbol); err == nil {
		if err = json.Unmarshal([]byte(coinBytes), &coin); err == nil {
			return coin, nil
		}
	}
	return dimMod.DimCoin{}, err
}

func init() {
	coinRows := make([]dimMod.DimCoin, 0)
	if err := models.GetBeans(&coinRows, dimMod.DimCoin{State: utils.STATE_ENABLE}); err == nil {
		for _, coin := range coinRows {
			// coins[coin.Symbol] = coin
			coinByts, _ := json.Marshal(coin)
			utils.Set("dimcoin:"+coin.Symbol, coinByts, 0)
		}
	}

}
func (DimCoinService) Find() ([]dimMod.DimCoin, error) {
	coins := make([]dimMod.DimCoin, 0)
	err := models.GetBeans(&coins, dimMod.DimCoin{State: utils.STATE_ENABLE})
	return coins, err
}

func (DimCoinService) GetExchangeRate(mainSymbol, exchangeSymbol string) (float64, error) {
	var mainCoin dimMod.DimCoin
	var exchangeCoin dimMod.DimCoin
	var mainCoinBytes, exchangeCoinBytes string
	var err error
	if mainCoinBytes, err = utils.Get("dimcoin:" + mainSymbol); err == nil {
		if exchangeCoinBytes, err = utils.Get("dimcoin:" + exchangeSymbol); err == nil {
			if err = json.Unmarshal([]byte(mainCoinBytes), &mainCoin); err == nil {
				if err = json.Unmarshal([]byte(exchangeCoinBytes), &exchangeCoin); err == nil {

					if mainCoin.PriceCny > 0 && exchangeCoin.PriceCny > 0 {
						return mainCoin.PriceCny / exchangeCoin.PriceCny, nil
					}
				}
			}

		}
	}
	return 0, fmt.Errorf("未找到[%s->%s]的汇率", mainCoin, exchangeCoin)
}

func (dimCoin DimCoinService) GetExchangeFree(symbol string) (float64, bool) {
	if coin, err := dimCoin.GetCoin(symbol); err == nil {
		return coin.ExchangeFree, true
	}
	return 0, false
}
func (dimCoin DimCoinService) GetFree(symbol string) (float64, bool) {
	if coin, err := dimCoin.GetCoin(symbol); err == nil {
		return coin.TransferFree, true
	}
	return 0, false
}
