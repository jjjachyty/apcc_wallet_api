package dappSrv

import (
	"apcc_wallet_api/models"
	"apcc_wallet_api/models/dappMod"
	"apcc_wallet_api/utils"
)

type DappUseLogService struct{}

func (DappUseLogService) RecentUsed(user string) ([]dappMod.Dapp, error) {
	dapp := make([]dappMod.Dapp, 0)
	session, _ := utils.GetSession()
	err := session.SQL(`
	select * from dapp app where app.uuid in  (
		SELECT  log.dapp
		FROM
			dapp_use_log log 
		WHERE
			log.user = ?
		GROUP BY log.dapp
		order by max(use_at)
		)
		LIMIT 0,6`, user).Find(&dapp)
	return dapp, err
}

func (DappUseLogService) Add(log dappMod.DappUseLog) error {
	return models.Create(log)
}
