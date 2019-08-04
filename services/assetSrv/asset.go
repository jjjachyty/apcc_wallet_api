package assetSrv

import (
	"apcc_wallet_api/models"
	"apcc_wallet_api/models/assetMod"
	"apcc_wallet_api/services/dimSrv"
	"apcc_wallet_api/services/userSrv"
	"apcc_wallet_api/utils"
	"database/sql"
	"errors"

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

var userService userSrv.UserService

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

func (AssetService) CreateLog(assetsLog []assetMod.AssetLog) error {

	return models.Create(&assetsLog)
}
func (AssetService) Get(assets *assetMod.Asset) error {
	return models.SQLBean(assets, getAssetsSQL, assets.UUID)
}

func (AssetService) GetLogs(page *utils.PageData, condBean interface{}) error {
	return models.GetBeansPage(page, condBean)
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

//Coin2MHC 货币兑换MHC
func (AssetService) Coin2MHC(log assetMod.AssetLog) error {

	return utils.Session(func(session *xorm.Session) (err error) {
		var result sql.Result
		var rows int64
		if err = session.Begin(); err == nil {
			if result, err = session.Exec("UPDATE asset a set a.blance = a.blance-? where a.uuid = ? and a.address=? and a.symbol=? and a.blance > ? ", log.FromAmount, log.FromUser, log.FromAddress, log.FromCoin, log.FromAmount); err == nil {
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
func (AssetService) MHC2Coin(log assetMod.AssetLog) error {

	return utils.Session(func(session *xorm.Session) (err error) {
		var result sql.Result
		var rows int64
		if err = session.Begin(); err == nil {
			if result, err = session.Exec("UPDATE asset a set a.blance = a.blance+? where a.uuid = ? and a.address = ? and  a.symbol=? ", log.ToAmount, log.FromUser, log.ToAddress, log.ToCoin); err == nil {
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

func (AssetService) Send(log assetMod.AssetLog) error {

	return utils.Session(func(session *xorm.Session) (err error) {
		var result sql.Result
		var rows int64
		if err = session.Begin(); err == nil {
			if result, err = session.Exec("UPDATE asset a set a.blance = a.blance-? where a.uuid = ? and a.symbol=? and a.blance>? ", log.FromAmount, log.FromUser, log.FromCoin, log.FromAmount); err == nil {
				if rows, err = result.RowsAffected(); err == nil {
					if rows == 1 {
						if _, err = session.Exec("UPDATE asset a set a.blance = a.blance+? where a.address = ? and a.symbol=? ", log.ToAmount, log.ToAddress, log.FromCoin, log.ToAmount); err == nil {
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

func (AssetService) UpdateLogs(log assetMod.AssetLog) error {
	return models.UpdateBean(log, assetMod.Asset{UUID: log.UUID})
}
