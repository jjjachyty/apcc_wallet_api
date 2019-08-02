package assetSrv

import (
	"apcc_wallet_api/models"
	"apcc_wallet_api/models/assetMod"
	"apcc_wallet_api/utils"

	"github.com/ethereum/go-ethereum/ethclient"
)

var client *ethclient.Client
var waitExchanges []assetMod.Exchange

type ExchangeService struct{}

//Add 新增转换记录
func (ExchangeService) Add(exchange *assetMod.Exchange) error {
	return models.Create(exchange)
}
func (ExchangeService) GetExchanges(page *utils.PageData, condBean interface{}) error {
	return models.GetBeansPage(page, condBean)

}

func getWaitExchange() ([]assetMod.Exchange, error) {
	exchanges := make([]assetMod.Exchange, 0)
	err := models.SQLBeans(&exchanges, "SELECT * from exchange exg where exg.state=? order by create_at,amount desc", utils.STATE_DISABLE)
	return exchanges, err
}

func init() {
	var err error
	if client, err = ethclient.Dial("http://119.3.108.19:8110"); err != nil {
		panic("MHC服务器连接失败")
	}

	// go func() {
	// 	timer := time.NewTimer(time.Second * 10)
	// 	for {
	// 		select {
	// 		case <-timer.C:
	// 			timer.Reset(time.Second * 10)
	// 			Scan()
	// 		}
	// 	}
	// }()
}
