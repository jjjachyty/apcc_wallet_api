package commonSrv

import (
	"apcc_wallet_api/utils"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
)

var CoinPrice = make(map[string]float64)
var CoinFreee = make(map[string]float64)

var apiURL = "https://www.hqz.com/api/index/get_index_refresh/?type=&sortby=&sort=desc&filter=&p=1"

func init() {
	GetCoinPrice()
	GetCoinFree()
}

func GetCoinFree() {
	CoinFreee["MHC"] = 1.0
	CoinFreee["USDT"] = 5.0
}

func GetCoinPrice() {
	CoinPrice["MedicalHealthCoin"] = 1

	if resp, err := http.Get(apiURL); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {

			bodyStr := string(body)
			coins := gjson.Get(bodyStr, "coin").Get("data").Array()

			for _, coin := range coins {
				CoinPrice[coin.Get("coincode").String()] = coin.Get("price_cny").Float()
			}
		}

	}
	if len(CoinPrice) == 0 {
		utils.AppLog.Errorln("获取行情失败")
	}

}
func GetExchange(mainCoin, exchangeCoin string) ([]float64, error) {
	var exchangeRate = CoinPrice[exchangeCoin]
	var mainRate = CoinPrice[mainCoin]
	if exchangeRate > 0 && mainRate > 0 {
		return []float64{exchangeRate, mainRate}, nil
	}
	return nil, fmt.Errorf("未找到[%s->%s]的汇率", exchangeCoin, mainCoin)
}
