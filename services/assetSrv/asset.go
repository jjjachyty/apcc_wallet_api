package assetSrv

import (
	"apcc_wallet_api/models"
	"apcc_wallet_api/models/assetMod"
	"apcc_wallet_api/services/userSrv"
	"apcc_wallet_api/utils"

	"github.com/go-xorm/xorm"
)

type PayTYpe int

const (
	HMC_COIN_SYMBOL  = "MHC"
	USDT_COIN_SYMBOL = "USDT"

	PAY_TYPE_EXCHANGE PayTYpe = 1000
	PAY_TYPE_TRANSFER_INNER
	PAY_TYPE_TRANSFER_OUTER
	PAY_TYPE_TRANSFER_UNFREEZE
)

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

func (AssetService) Get(assets *assetMod.Asset) error {
	return models.SQLBean(assets, getAssetsSQL, assets.UUID)
}

func (AssetService) Find(uuid string) ([]assetMod.Asset, error) {
	var assets = make([]assetMod.Asset, 0)
	return assets, models.SQLBeans(&assets, assetsSQL, uuid)
}

func (AssetService) FindExchange(uuid string, mainCoin string, exchangeCoin string) ([]assetMod.Asset, error) {
	var assets = make([]assetMod.Asset, 0)
	return assets, models.SQLBeans(&assets, exchangeSQL, uuid, mainCoin, exchangeCoin)
}

func (AssetService) FindInnerTransfer(from, to assetMod.Asset) ([]assetMod.Asset, error) {
	var assets = make([]assetMod.Asset, 0)
	return assets, models.SQLBeans(&assets, innerTransferSQL, from.UUID, from.Address, to.Address, from.Symbol)
}

func (AssetService) ExchangeCoin(from, to assetMod.Asset, transferAmount float64) error {
	var exchangeRate = from.PriceCny / to.PriceCny
	var toAmount = exchangeRate * transferAmount
	return utils.Session(func(session *xorm.Session) (err error) {
		if err = session.Begin(); err == nil {
			if _, err = session.Exec("UPDATE asset a set a.blance = ? where a.uuid = ? and a.symbol=? and a.blance=? ", from.Blance-transferAmount, from.UUID, from.Symbol, from.Blance); err == nil {
				if _, err = session.Exec("UPDATE asset a set a.blance = ? where a.uuid = ? and a.symbol=? and a.blance=? ", to.Blance+toAmount, to.UUID, to.Symbol, to.Blance); err == nil {
					_, err = session.Insert(assetMod.AssetLog{UUID: utils.GetUUID(), FromUser: from.UUID, FromAddress: from.Address, FromPreblance: from.Blance, FromBlance: from.Blance - transferAmount, FromPriceCny: from.PriceCny,
						ToUser: to.UUID, ToAddress: to.Address, ToPreblance: to.Blance, ToBlance: to.Blance + toAmount, ToPriceCny: to.PriceCny, PayType: int(PAY_TYPE_EXCHANGE), State: utils.STATE_ENABLE})
				}
			}
			if err == nil {
				err = session.Commit()
			} else {
				err = session.Rollback()
			}
		}
		return err

	})

}

func (AssetService) Send(from, to assetMod.Asset, amount float64, payType PayTYpe) error {

	return utils.Session(func(session *xorm.Session) (err error) {
		var state = utils.STATE_ENABLE
		if payType == PAY_TYPE_TRANSFER_OUTER {
			state = utils.STATE_DISABLE
		}
		if err = session.Begin(); err == nil {
			if _, err = session.Exec("UPDATE asset a set a.blance = ? where a.uuid = ? and a.symbol=? and a.blance=? ", from.Blance-amount, from.UUID, from.Symbol, from.Blance); err == nil {
				if _, err = session.Exec("UPDATE asset a set a.blance = ? where a.uuid = ? and a.symbol=? and a.blance=? ", to.Blance+amount, to.UUID, to.Symbol, to.Blance); err == nil {
					_, err = session.Insert(assetMod.AssetLog{UUID: utils.GetUUID(), FromUser: from.UUID, FromAddress: from.Address, FromPreblance: from.Blance, FromBlance: from.Blance - amount, FromPriceCny: from.PriceCny,
						ToUser: to.UUID, ToAddress: to.Address, ToPreblance: to.Blance, ToBlance: to.Blance + amount, ToPriceCny: from.PriceCny, PayType: int(payType), State: state})
				}
			}
			if err == nil {
				err = session.Commit()
			} else {
				err = session.Rollback()
			}
		}
		return err

	})

}
