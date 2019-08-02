package dappMod

import "time"

//Dapp Dapp
type DappUseLog struct {
	UUID     string `xorm:"varchar(36) 'uuid'"`
	User     string
	Dapp     string
	UseAt    time.Time `xorm:"created"`
	Name     string    `xorm:"<-`
	Logo     string    `xorm:"<-`
	HomePage string    `xorm:"<-`
}
