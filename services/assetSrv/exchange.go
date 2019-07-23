package assetSrv

import (
	"apcc_wallet_api/models"
	"apcc_wallet_api/models/assetMod"
)

type ExchangeService struct{}

//Add 新增转换记录
func (ExchangeService) Add(exchange *assetMod.Exchange) error {
	return models.Create(exchange)
}
