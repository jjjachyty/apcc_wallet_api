package userSrv

import (
	"apcc_wallet_api/models"
	"apcc_wallet_api/models/userMod"
	"apcc_wallet_api/utils"
)

type UserService struct{}

func (UserService) Register(user *userMod.User) error {
	user.Password = utils.GetMD5(user.Password)
	user.NickName = user.Phone[7:]
	return models.Create(user)
}

func (UserService) Login(user *userMod.User) error {
	user.Password = utils.GetMD5(user.Password)

	return models.GetBean(user)
}

func (UserService) Get(user *userMod.User) error {

	return models.GetBean(user)

}

func (UserService) Update(user *userMod.User) error {

	return models.UpdateBean(user, userMod.User{Phone: user.Phone})

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
