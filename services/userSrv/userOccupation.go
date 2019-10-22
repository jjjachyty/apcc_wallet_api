package userSrv

import (
	"apcc_wallet_api/models"
	"apcc_wallet_api/models/userMod"
)

type UserOccupationService struct {
}

func (UserOccupationService) Add(userOcc *userMod.UserOccupation) error {
	return models.Create(userOcc)
}
