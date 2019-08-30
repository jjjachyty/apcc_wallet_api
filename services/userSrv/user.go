package userSrv

import (
	"apcc_wallet_api/models"
	"apcc_wallet_api/models/assetMod"
	"apcc_wallet_api/models/userMod"
	"apcc_wallet_api/services/walletSrv"
	"apcc_wallet_api/utils"

	"github.com/go-xorm/xorm"
)

type UserService struct{}

func (UserService) Register(user *userMod.User) (err error) {
	var assets []assetMod.Asset

	return utils.Session(func(session *xorm.Session) error {
		user.UUID = utils.GetUUID()
		user.State = utils.STATE_ENABLE
		user.Password = utils.GetMD5(user.Password)
		user.NickName = user.Phone[7:]
		if _, err = session.Insert(user); err == nil {

			if assets, err = walletSrv.GetAddress(user.UUID, uint32(user.AccountID)); err == nil {
				_, err = session.Insert(assets)
			}

		}
		if err != nil {
			session.Rollback()
		} else {
			session.Commit()
		}
		return err
	})
}

func (UserService) Login(user *userMod.User) error {
	user.Password = utils.GetMD5(user.Password)

	return models.GetBean(user)
}

func (UserService) Get(user *userMod.User) error {

	return models.GetBean(user)

}

func (UserService) Update(user *userMod.User) error {
	uuid := user.UUID
	user.UUID = ""
	return models.UpdateBean(user, userMod.User{UUID: uuid})

}

func (this UserService) CheckPayPasswd(user *userMod.User) (bool, error) {
	var err error
	if err = this.Get(user); err == nil {
		if user.Phone != "" {
			return true, nil
		}
	}
	return false, err
}

func (UserService) UpdateIDCard(card *userMod.IdCard) error {

	return models.UpdateBean(card, userMod.IdCard{UserID: card.UserID})

}

func (UserService) GetMaxCountID() int64 {
	var maxCountID int64
	models.SQLBean(&maxCountID, "SELECT MAX(account_id) FROM user ")
	return maxCountID
}
