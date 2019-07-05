package assetMod

import (
	"apcc_wallet_api/models/dimMod"
	"time"
)

type Asset struct {
	UUID           string `xorm:"varchar(36) notnull unique pk 'uuid'"`
	Symbol         string
	Address        string
	Blance         float64 `json:"Blance,string"`
	FreezingBlance float64 `json:"FreezingBlance,string"`
	dimMod.DimCoin `xorm:"extends"`
}

func (Asset) TableName() string {
	return "asset"
}

type AssetLog struct {
	UUID          string `xorm:"varchar(36) 'uuid'"`
	FromAddress   string
	FromUser      string
	FromPreblance float64
	FromBlance    float64
	FromPriceCny  float64
	ToUser        string
	ToAddress     string
	ToPreblance   float64
	ToBlance      float64
	ToPriceCny    float64
	CreateAt      time.Time `xorm:"created"`
	PayType       int
	State         int
}

func (AssetLog) TableName() string {
	return "asset_log"
}
