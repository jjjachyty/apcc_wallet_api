package assetMod

import "apcc_wallet_api/models/dimMod"

type Asset struct {
	UUID           string `xorm:"varchar(36) notnull unique pk 'uuid'"`
	Symbol         string
	Address        string
	Blance         string
	FreezingBlance string
	PriceCny       float64 `xorm:"-"`
	dimMod.DimCoin
}
