package assetSrv

import (
	"apcc_wallet_api/models"
	"apcc_wallet_api/models/assetMod"
	"apcc_wallet_api/utils"
	"database/sql"
	"errors"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-xorm/xorm"
)

var client *ethclient.Client
var waitExchanges []assetMod.ExchangeLog

type ExchangeService struct{}

//Add 新增转换记录
func (ExchangeService) Add(exchange *assetMod.ExchangeLog) error {
	return models.Create(exchange)
}
func (ExchangeService) GetExchanges(page *utils.PageData, condBean interface{}) error {
	return models.GetBeansPage(page, condBean)
}
func (ExchangeService) Update(exchnage *assetMod.ExchangeLog) error {
	return models.Update(exchnage.UUID, exchnage)
}

//Coin2MHC 货币兑换MHC
func (ExchangeService) Coin2MHC(log assetMod.ExchangeLog) error {

	return utils.Session(func(session *xorm.Session) (err error) {
		var result sql.Result
		var rows int64
		if err = session.Begin(); err == nil {
			if result, err = session.Exec("UPDATE asset a set a.blance = a.blance-? where a.uuid = ? and a.address=? and a.symbol=? and a.blance > ? ", log.FromAmount, log.User, log.FromAddress, log.FromCoin, log.FromAmount); err == nil {
				if rows, err = result.RowsAffected(); err == nil {
					if rows == 1 {
						_, err = session.Insert(log)
					} else {
						err = errors.New("兑换失败,请检查是否有足够的余额")
					}
				}
			}
			if err == nil {
				err = session.Commit()
			} else {
				session.Rollback()
			}
		}
		return err

	})

}

//MHC2Coin MHC兑换货币
func (ExchangeService) MHC2Coin(log assetMod.ExchangeLog) error {

	return utils.Session(func(session *xorm.Session) (err error) {
		var result sql.Result
		var rows int64
		if err = session.Begin(); err == nil {
			if result, err = session.Exec("UPDATE asset a set a.blance = a.blance+? where a.uuid = ? and a.address = ? and  a.symbol=? ", log.ToAmount, log.User, log.ToAddress, log.ToCoin); err == nil {
				if rows, err = result.RowsAffected(); err == nil {
					if rows == 1 {
						_, err = session.Insert(log)
					} else {
						err = errors.New("兑换失败,请检查账户是否正常")
					}
				}
			}
			if err == nil {
				err = session.Commit()
			} else {
				session.Rollback()
			}
		}
		return err

	})

}
