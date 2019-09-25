package imCtr

import (
	"apcc_wallet_api/middlewares/jwt"
	"apcc_wallet_api/services/imSrv"
	"apcc_wallet_api/utils"

	"github.com/gin-gonic/gin"
)

type FriendController struct{}

func (FriendController) Apply(c *gin.Context) {
	var err error
	var friendApply = new(imSrv.FriendApply)
	if err = c.BindJSON(friendApply); err == nil {
		claims := jwt.GetClaims(c)
		friendApply.UserID = claims.UUID
		err = friendApply.Create()
	}
	utils.Response(c, err, nil)
}

func (FriendController) Agree(c *gin.Context) {
	var err error
	var friendApply = new(imSrv.FriendApply)
	if err = c.BindJSON(friendApply); err == nil {
		err = friendApply.Agree()
	}
	utils.Response(c, err, nil)
}

func (FriendController) Del(c *gin.Context) {
	var err error
	var user = new(imSrv.Friend)
	if err = c.BindJSON(user); err == nil {
		err = user.Delete()
	}
	utils.Response(c, err, nil)
}
func (FriendController) Update(c *gin.Context) {
	var err error
	var user = new(imSrv.Friend)
	if err = c.BindJSON(user); err == nil {
		err = user.Update()
	}
	utils.Response(c, err, nil)

}
func (FriendController) GetApply(c *gin.Context) {
	var err error
	claims := jwt.GetClaims(c)
	var applys []imSrv.FriendApply
	applys, err = imSrv.FriendApplyService{}.Get(claims.UUID)
	utils.Response(c, err, applys)
}

func (FriendController) UpdateApply(c *gin.Context) {
	var err error
	var friendsApply = new(imSrv.FriendApply)
	if err = c.BindJSON(friendsApply); err == nil {
		err = friendsApply.Update()
	}
	utils.Response(c, err, friendsApply)
}

func (FriendController) GetFriends(c *gin.Context) {
	var err error
	var friends []imSrv.Friend
	claims := jwt.GetClaims(c)
	friends, err = imSrv.FriendService{}.Get(claims.UUID)

	utils.Response(c, err, friends)
}
