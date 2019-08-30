package assetSrv

import (
	"apcc_wallet_api/models"
	"apcc_wallet_api/models/assetMod"
	"apcc_wallet_api/services/dimSrv"
	"apcc_wallet_api/utils"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-xorm/xorm"
)

type PayTYpe int

const (
	HMC_COIN_SYMBOL  = "MHC"
	USDT_COIN_SYMBOL = "USDT"
)
const (
	PAY_TYPE_EXCHANGE int = iota + 1000
	PAY_TYPE_TRANSFER_INNER
	PAY_TYPE_TRANSFER_OUTER
	PAY_TYPE_TRANSFER_UNFREEZE
	PAY_TYPE_TRANSFER_ADD_IN
)

var dimCoinService dimSrv.DimCoinService
var (
	getAssetsSQL     = "select * from asset where uuid = ?"
	assetsSQL        = "select a.*,b.name_en ,b.name_cn,b.price_cny,b.price_usd FROM asset a LEFT JOIN dim_coin b on  a.symbol = b.symbol where a.uuid=?"
	exchangeSQL      = "select a.*,b.name_en ,b.name_cn,b.price_cny,b.price_usd FROM asset a LEFT JOIN dim_coin b on  a.symbol = b.symbol where a.uuid=? and (a.symbol=? Or a.symbol = ?)"
	innerTransferSQL = `SELECT
								a.*, b.name_en,
								b.name_cn,
								b.price_cny,
								b.price_usd
						FROM
									asset a
						LEFT JOIN dim_coin b ON a.symbol = b.symbol
						WHERE
									(a.uuid =? AND a.address =?)
						UNION
						SELECT
									a.*, b.name_en,
									b.name_cn,
									b.price_cny,
									b.price_usd
						FROM
									asset a
						LEFT JOIN dim_coin b ON a.symbol = b.symbol
						WHERE
									( a.address=? and a.symbol=? )`
)

type AssetService struct{}

func TransferCoin(symbol string, fromAddress string, amount float64, toAddress string) {
	switch symbol {
	case HMC_COIN_SYMBOL:

	case USDT_COIN_SYMBOL:
	}
}

func (AssetService) Create(assets []assetMod.Asset) error {

	return models.Create(&assets)
}
func (AssetService) Update(assets []assetMod.Asset) error {
	return models.Create(&assets)
}

func (AssetService) CreateLog(assetsLog assetMod.TransferLog) error {

	return models.Create(&assetsLog)
}
func (AssetService) Get(assets *assetMod.Asset) error {
	return models.SQLBean(assets, getAssetsSQL, assets.UUID)
}
func (AssetService) GetBean(assets *assetMod.Asset) error {
	return models.GetBean(assets)
}
func (AssetService) GetLogs(page *utils.PageData, log assetMod.TransferLog) error {
	return models.GetSQLPage(page, "coin = ? and (from_address = ? or to_address = ?)", log.Coin, log.FromAddress, log.ToAddress)
}

func (AssetService) Find(assets *[]assetMod.Asset, condBean assetMod.Asset) error {
	return models.GetBeans(assets, condBean)
}

func (AssetService) FindExchange(uuid string, mainCoin string, exchangeCoin string) ([]assetMod.Asset, error) {
	var assets = make([]assetMod.Asset, 0)
	return assets, models.SQLBeans(&assets, exchangeSQL, uuid, mainCoin, exchangeCoin)
}

func (AssetService) FindInnerTransfer(from, to assetMod.Asset) ([]assetMod.Asset, error) {
	var assets = make([]assetMod.Asset, 0)
	return assets, models.SQLBeans(&assets, innerTransferSQL, from.UUID, from.Address, to.Address, from.Symbol)
}

func (AssetService) SendInner(log assetMod.TransferLog) error {

	return utils.Session(func(session *xorm.Session) (err error) {
		var result sql.Result
		var rows int64
		if err = session.Begin(); err == nil {
			if result, err = session.Exec("UPDATE asset a set a.blance = a.blance-? where a.uuid = ? and a.address=? and a.symbol=? and a.blance>=? ", log.Amount-log.Free, log.FromUser, log.FromAddress, log.Coin, log.Amount-log.Free); err == nil {
				if rows, err = result.RowsAffected(); err == nil {
					if rows == 1 {
						if _, err = session.Exec("UPDATE asset a set a.blance = a.blance+? where a.address = ? and a.symbol=? ", log.Amount-log.Free, log.ToAddress, log.Coin); err == nil {
							_, err = session.Insert(log)
						}

					} else {
						err = errors.New("转账失败,请检查是否有足够的余额")
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

func (AssetService) SendOuter(log assetMod.TransferLog) error {

	return utils.Session(func(session *xorm.Session) (err error) {
		var result sql.Result
		var rows int64
		if err = session.Begin(); err == nil {
			if result, err = session.Exec("UPDATE asset a set a.blance = a.blance-? where a.uuid = ? and a.address=? and a.symbol=? and a.blance>=? ", log.Amount, log.FromUser, log.FromAddress, log.Coin, log.Amount); err == nil {
				if rows, err = result.RowsAffected(); err == nil {
					if rows == 1 {
						_, err = session.Insert(log)
					} else {
						err = errors.New("转账失败,请检查是否有足够的余额")
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
func (AssetService) UpdateTransferLog(log assetMod.TransferLog) error {
	return models.UpdateBean(log, assetMod.Asset{UUID: log.UUID})
}

//AddCoin Coin转入平台
func (AssetService) AddCoin(log assetMod.TransferLog) error {
	return utils.Session(func(session *xorm.Session) (err error) {
		var result sql.Result
		var rows int64
		if err = session.Begin(); err == nil {
			if result, err = session.Exec("UPDATE asset a set a.blance = a.blance+? where  a.address = ? and  a.symbol=? ", log.Amount, log.ToAddress, log.Coin); err == nil {
				if rows, err = result.RowsAffected(); err == nil {
					if rows == 1 {
						_, err = session.Update(assetMod.TransferLog{State: utils.STATE_ENABLE}, assetMod.TransferLog{UUID: log.UUID})
					} else {
						err = fmt.Errorf("入账失败,请检查账户%s 币种 %s 是否正常", log.ToAddress, log.Coin)
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
