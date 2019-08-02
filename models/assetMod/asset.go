package assetMod

import (
	"time"
)

type Asset struct {
	UUID           string `xorm:"varchar(36) notnull unique pk 'uuid'"`
	Symbol         string
	BaseOn         string
	Address        string
	Blance         float64
	FreezingBlance float64
	// dimMod.DimCoin `xorm:"<-"`
}

func (Asset) TableName() string {
	return "asset"
}

type AssetLog struct {
	UUID        string `xorm:"varchar(36) 'uuid'"`
	FromAddress string
	FromUser    string
	FromCoin    string

	FromPriceCny float64
	// ExchangeAddress string
	ExchangeTxs string
	ToUser      string
	ToCoin      string
	ToAddress   string

	ToPriceCny float64
	CreateAt   time.Time `xorm:"created"`
	FromAmount float64
	ToAmount   float64

	Free        float64
	PayType     int
	State       int
	SendTxs     string
	SendAddress string
	SendTime    time.Time `xorm:"updated"`
}

func (AssetLog) TableName() string {
	return "asset_log"
}
