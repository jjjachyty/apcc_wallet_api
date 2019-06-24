package authService

import (
	"apcc_wallet_api/models"
	"apcc_wallet_api/models/authModel"
	"apcc_wallet_api/utils"
)

type UserService struct{}

func (UserService) Register(user *authModel.User) error {
	user.Password = utils.GetMD5(user.Password)
	user.NickName = user.Phone[7:]
	return models.Create(user)
}

func (UserService) Login(user *authModel.User) error {
	user.Password = utils.GetMD5(user.Password)

	return models.GetBean(user)
}

func (UserService) Get(user *authModel.User) error {

	return models.GetBean(user)

}
