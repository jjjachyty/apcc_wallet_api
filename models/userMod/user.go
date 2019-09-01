package userMod

import (
	"fmt"
	"time"
)

type Birthday time.Time

var birthdayFormart = "2006-01-02"

func (t Birthday) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(birthdayFormart)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, birthdayFormart)
	b = append(b, '"')
	return b, nil
}

func (t Birthday) String() string {
	return time.Time(t).Format(birthdayFormart)
}

func (birthday *Birthday) UnmarshalJSON(b []byte) error {
	fmt.Println("time", b, string(b)[1:11])
	date, err := time.ParseInLocation(birthdayFormart, string(b)[1:11], time.Local)
	fmt.Println("date", date)
	*birthday = Birthday(date)
	return err
}

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
	UserID       string   `xorm:" notnull unique pk 'user_id'"`
	Gender       string   `xorm:"'gender'"`
	Name         string   `xorm:"'name'"`
	IDCardNumber string   `json:"id_card_number" xorm:"'number'"`
	Birthday     Birthday `xorm:"'birthday'"`
	Race         string   `xorm:"'race'"`
	Address      string   `xorm:"'address'"`
	// Legality": {
	// 	"Edited": 0.001,
	// 	"Photocopy": 0.0,
	// 	"ID Photo": 0.502,
	// 	"Screen": 0.496,
	// 	"Temporary ID Photo": 0.0
	// },
	// Type      int
	// Side      string
	IssuedBy  string `json:"issued_by" xorm:"'issued_by'"`
	ValidDate string `json:"valid_date" xorm:"'valid_date'"`
}
