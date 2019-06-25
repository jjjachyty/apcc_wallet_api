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
