package userMod

import (
	"apcc_wallet_api/models/coinMod"
	"time"
)

type User struct {
	UUID            string `xorm:"varchar(26) notnull unique pk 'uuid'"`
	Phone           string `xorm:"varchar(25) notnull unique pk 'phone'"`
	NickName        string
	Avatar          string
	Password        string
	HasPayPasswd    bool `xorm:"-"` //是否有交易密码
	PayPasswd       string
	LastLoginTime   time.Time
	LastLoginIP     string `xorm:"varchar(25) 'last_login_ip'"`
	LastLoginDevice string
	Identification  int             // 账户标识/新客户、老客户
	Level           int             //账户等级
	State           string          //账户状态
	IDCardAuth      bool            `xorm:"-"`
	IDCard          Card            `xorm:"extends"`
	Accounts        coinMod.Account `xorm:"-"` //账户
}

type Card struct {
	Gender       string `xorm:"'id_gender'"`
	Name         string `xorm:"'id_name'"`
	IDCardNumber string `json:"id_card_number" xorm:"'id_card_number'"`
	Birthday     string `xorm:"'id_birthday'"`
	Race         string `xorm:"'id_race'"`
	Address      string `xorm:"'id_address'"`
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
