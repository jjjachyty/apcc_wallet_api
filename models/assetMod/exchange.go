package assetMod

import "time"

//Exchange 币种兑换表
type ExchangeLog struct {
	UUID         string `xorm:"varchar(36)  notnull unique pk  'uuid'"`
	User         string
	FromCoin     string
	FromAddress  string
	FromPriceCny float64
	FromAmount   float64

	ToCoin     string
	ToAddress  string
	ToPriceCny float64
	ToAmount   float64

	SendAddress string
	SendTxs     string
	SendAt      time.Time
	Free        float64
	CreateAt    time.Time `xorm:"created"`
	State       int
}
