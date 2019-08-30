package dappMod

import "time"

//Dapp Dapp
type Dapp struct {
	UUID          string `xorm:"varchar(36) 'uuid'"`
	Name          string
	Subtitle      string
	Category      string
	Permission    string
	Synopsis      string
	Score         float64
	Logo          string
	Banner        string
	FullScreen    int
	Snapshot      string
	Video         string
	Owner         string
	Used          string
	State         int
	SubmitAt      time.Time `xorm:"created"`
	UpperAt       time.Time
	LowerAt       time.Time
	Auditor       string
	AuditOpinions string
	AbiCode       string
	Address       string
	HomePage      string
}

// DappSearchSQL Dapp 查询SQL
var DappSearchSQL = map[string]string{
	"Name": "name like ?",
}
