package commonSrv

import (
	"apcc_wallet_api/utils"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
)

var CoinPrice = make(map[string]float64)

var apiURL = "https://www.hqz.com/api/index/get_index_refresh/?type=&sortby=&sort=desc&filter=&p=1"

func init() {
	if resp, err := http.Get(apiURL); err == nil {
		defer resp.Body.Close()
		if body, err := ioutil.ReadAll(resp.Body); err == nil {
			bodyStr := string(body)
			coins := gjson.Get(bodyStr, "data").Array()
			for _, coin := range coins {
				CoinPrice[coin.Get("coincode").String()] = coin.Get("price_cny").Float()
			}
		}

	}
	if len(CoinPrice) == 0 {
		utils.AppLog.Errorln("获取行情失败")
	}
}
