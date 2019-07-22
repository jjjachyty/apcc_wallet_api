package dappSrv

import (
	"apcc_wallet_api/models"
	"apcc_wallet_api/models/dappMod"
	"apcc_wallet_api/utils"
	"fmt"
	"net/url"
)

type DappService struct{}

func (DappService) Page(page *utils.PageData, pars url.Values) error {
	sql, args := models.GetSQL(pars, dappMod.DappSearchSQL)
	fmt.Println("SLQ=====", sql, args)
	return models.GetSQLPage(page, sql, args...)
}
func (DappService) Find(results interface{}, condBean interface{}) error {
	return models.GetBeans(results, condBean)
}
