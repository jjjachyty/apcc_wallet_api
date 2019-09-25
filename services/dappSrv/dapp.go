package dappSrv

import (
	"apcc_wallet_api/models"
	"apcc_wallet_api/models/dappMod"
	"apcc_wallet_api/utils"
	"net/url"
)

type DappService struct{}

func (DappService) Page(page *utils.PageData, pars url.Values) error {
	sql, args := models.GetSQL(pars, dappMod.DappSearchSQL)
	return models.GetSQLPage(page, sql, args...)
}
func (DappService) Find(results interface{}, condBean interface{}) error {
	return models.GetBeans(results, condBean)
}

func (DappService) AdvancedFind(orderDesc string, limit int, where string, whrerargs ...interface{}) ([]dappMod.Dapp, error) {
	dapps := make([]dappMod.Dapp, 0)
	session, _ := utils.GetSession()
	err := session.Desc(orderDesc).Limit(limit).Where(where, whrerargs...).Find(&dapps)
	return dapps, err
}
