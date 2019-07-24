package assetMod

import "time"

//Exchange 币种兑换表
type Exchange struct {
	UUID           string `xorm:"varchar(36) 'uuid'"`
	User           string
	FromCoin       string
	FromAddress    string
	ReceiveAddress string
	ReceiveTxs     string
	ToCoin         string
	ToAddress      string
	SendTxs        string
	SendAt         time.Time
	Amount         float64
	Free           float64
	Rate           float64
	CreateAt       time.Time `xorm:"created"`
	State          int
}
