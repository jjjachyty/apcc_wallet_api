package assetSrv

import (
	"apcc_wallet_api/models"
	"apcc_wallet_api/models/assetMod"
	"apcc_wallet_api/utils"
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

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
	Scan()
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

func Scan() {
	if exchanges, err := getWaitExchange(); err == nil {

		for _, exchange := range exchanges {
			switch exchange.FromCoin + exchange.ToCoin {
			case "MHCUSDT":
				mhc2usdt(exchange)
			case "USDTMHC":
			}
		}
	}
}

//mhc 兑换usdt
func mhc2usdt(exchg assetMod.Exchange) {
	if tx, isPending, err := client.TransactionByHash(context.Background(), common.HexToHash(exchg.ReceiveTxs)); err == nil && !isPending {
		amount := tx.Value()
		fmt.Println("到账 amount", amount.String())

		tenDecimal := big.NewFloat(math.Pow(10, float64(18)))
		freeAmount, _ := new(big.Float).Mul(tenDecimal, big.NewFloat(exchg.Free)).Int(&big.Int{})
		fmt.Println("手续费", freeAmount.String())

		rateAmount, _ := big.NewFloat(exchg.Rate * 1000000).Int(&big.Int{})
		fmt.Println("费率", rateAmount.String())

		fmt.Printf("原始金额 amount %f \n", exchg.Amount)

		transAmount := amount.Sub(amount, freeAmount)

		ustdAmount := new(big.Int).Mul(transAmount, rateAmount)
		fmt.Println("ustdAmount====", ustdAmount.String())
		divUnit, _ := big.NewFloat(0).SetString("1000000000000000000000000")
		transAmountf := new(big.Float).Quo(big.NewFloat(0).SetInt(ustdAmount), divUnit)

		fmt.Println("USDT AMount")
		fmt.Println(transAmountf.Float64())

	}

}
