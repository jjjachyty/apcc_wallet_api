package userMod

import (
	"time"
)

type User struct {
	UUID         string `xorm:"varchar(26) notnull unique pk 'uuid'"`
	AccountID    int    `xorm:"int(11)  autoincr 'account_id'"`
	Phone        string `xorm:"varchar(25) notnull unique pk 'phone'"`
	NickName     string
	Avatar       string
	Password     string
	HasPayPasswd bool `xorm:"-"` //是否有交易密码
	PayPasswd    string
	// CreateAt        time.Time
	// LastLoginTime   time.Time
	LastLoginIP     string `xorm:"varchar(25) 'last_login_ip'"`
	LastLoginDevice string
	Identification  int // 账户标识/新客户、老客户
	Level           int //账户等级
	State           int //账户状态
	IDCardAuth      int `xorm:"'id_card_auth'"`
}

type IdCard struct {
	UserID       string    `xorm:" notnull unique pk 'user_id'"`
	Gender       string    `xorm:"'id_gender'"`
	Name         string    `xorm:"'id_name'"`
	IDCardNumber string    `json:"id_card_number" xorm:"'id_card_number'"`
	Birthday     time.Time `xorm:"'id_birthday'"`
	Race         string    `xorm:"'id_race'"`
	Address      string    `xorm:"'id_address'"`
	// Legality": {
	// 	"Edited": 0.001,
	// 	"Photocopy": 0.0,
	// 	"ID Photo": 0.502,
	// 	"Screen": 0.496,
	// 	"Temporary ID Photo": 0.0
	// },
	// Type      int
	// Side      string
	IssuedBy  string `json:"issued_by" xorm:"'id_issued_by'"`
	ValidDate string `json:"valid_date" xorm:"'id_valid_date'"`
}
