package userCtr

import (
	"apcc_wallet_api/middlewares/jwt"
	"apcc_wallet_api/models/userMod"
	"apcc_wallet_api/services/userSrv"
	"apcc_wallet_api/utils"

	"github.com/gin-gonic/gin"
)

type UserOccupationController struct{}

var userOccupationService userSrv.UserOccupationService

func (UserOccupationController) Apply(c *gin.Context) {
	var err error
	var user = new(userMod.UserOccupation)
	if err = c.BindJSON(user); err == nil {
		claims := jwt.GetClaims(c)
		user.UserID = claims.UUID
		err = userOccupationService.Add(user)

	}
	utils.Response(c, err, nil)

}
