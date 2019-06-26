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
	IDCard          IDCard          `xorm:"extends"`
	Accounts        coinMod.Account `xorm:"-"` //账户
}

type IDCard struct {
	IDCard         string    `xorm:"varchar(25) 'id_card'"` //身份证号
	Sex            string    //性别
	Birthday       time.Time //生日
	ExpirationDate time.Time //失效日期
}
