package imSrv

import (
	"apcc_wallet_api/models"
	"apcc_wallet_api/utils"
	"time"
)

type FriendApplyService struct{}

type FriendApply struct {
	UUID            string `xorm:"varchar(26) notnull unique pk 'uuid'"`
	UserID          string `xorm:"varchar(26) notnull  'user_id'"`
	FriendID        string `xorm:"varchar(26) notnull  'friend_id'"`
	Explain         string
	Reply           string
	CreateTime      time.Time `xorm:"created"`
	UpdateTime      time.Time `xorm:"updated"`
	FriendAvatar    string    `xorm:"<- 'avatar'"`
	FriendNickName  string    `xorm:"<- 'nick_name'"`
	FriendPhone     string    `xorm:"<- 'phone'"`
	FriendIntroduce string    `xorm:"<- 'introduce'"`
}

func (FriendApply) TableName() string {
	return "im_friend_apply"
}

//Create 添加好友请求
func (fd *FriendApply) Create() error {
	return models.Create(fd)
}

func (fd *FriendApply) Agree() error {
	var err error
	se, _ := utils.GetSession()

	if err = se.Begin(); err == nil {
		fds := make([]Friend, 2)
		fds[0] = Friend{UUID: utils.GetUUID(), UserID: fd.UserID, FriendID: fd.FriendID}
		fds[1] = Friend{UUID: utils.GetUUID(), UserID: fd.FriendID, FriendID: fd.UserID}
		err = models.Create(fds)
		if err == nil {
			if err = fd.Delete(); err == nil {
				se.Commit()
			}
		}
		if err != nil {
			se.Rollback()
		}
	}
	return err
}

//Delete 删除好友请求
func (fd *FriendApply) Delete() error {
	return models.Delete(fd.UUID, FriendApply{})
}

//Update 更新好友请求
func (fd *FriendApply) Update() error {
	return models.Update(fd.UUID, fd, []string{"explain", "reply", "update_time"}...)
}

//Get 查询我的好友
func (FriendApplyService) Get(userID string) ([]FriendApply, error) {
	friendsApply := make([]FriendApply, 0)
	err := models.SQLBeans(&friendsApply, `SELECT
	*
FROM
im_friend_apply friend
LEFT JOIN user u ON u.uuid = friend.user_id
WHERE
	friend_id = ? order by update_time desc`, userID)

	return friendsApply, err
}
