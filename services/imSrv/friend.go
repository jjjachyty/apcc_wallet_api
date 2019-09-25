package imSrv

import (
	"apcc_wallet_api/models"
	"time"
)

type FriendService struct{}

type Friend struct {
	UUID            string `xorm:"varchar(26) notnull unique pk 'uuid'"`
	UserID          string `xorm:"varchar(26) notnull  'user_id'"`
	FriendID        string `xorm:"varchar(26) notnull  'friend_id'"`
	Label           string
	CreateTime      time.Time `xorm:"created"`
	UpdateTime      time.Time `xorm:"updated"`
	FriendAvatar    string    `xorm:"<- 'avatar'"`
	FriendNickName  string    `xorm:"<- 'nick_name'"`
	FriendPhone     string    `xorm:"<- 'phone'"`
	FriendIntroduce string    `xorm:"<- 'introduce'"`
}

//Get 查询我的好友
func (FriendService) Get(userID string) ([]Friend, error) {
	friends := make([]Friend, 0)
	return friends, models.SQLBeans(&friends, `SELECT
	*
FROM
	im_friend friend
LEFT JOIN user u ON u.uuid = friend.friend_id
WHERE
	user_id = ?`, userID)
}

func (Friend) TableName() string {
	return "im_friend"
}

//Delete 删除好友
func (fd *Friend) Delete() error {
	return models.Delete(fd.UUID, Friend{})
}

//Update 更新好友
func (fd *Friend) Update() error {
	return models.Update(fd.UUID, fd, []string{"label", "update_time"}...)
}
