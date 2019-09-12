package assetMod

import (
	"time"
)

type Asset struct {
	UUID           string `xorm:"varchar(36) notnull unique pk 'uuid'"`
	Symbol         string
	BaseOn         string
	Address        string
	Blance         float64
	FreezingBlance float64
	// dimMod.DimCoin `xorm:"<-"`
}

func (Asset) TableName() string {
	return "asset"
}

type TransferLog struct {
	UUID        string `xorm:"varchar(36)  notnull unique pk  'uuid'"`
	FromAddress string
	FromUser    string
	Coin        string

	ToUser string

	ToAddress string

	PriceCny float64
	CreateAt time.Time `xorm:"created"`
	Amount   float64

	Free        float64
	PayType     int
	State       int
	SendTxs     string
	SendAddress string
	SendAt      time.Time
}

func (TransferLog) TableName() string {
	return "transfer_log"
}

type MHCTransferLog struct {
	TxHash      string `xorm:"varchar(36)  notnull unique pk"`
	BlockNumber int64
	BlockHash   string
	From        string
	To          string
	Gas         float64

	GasPrice float64
	GasUsed  float64
	Value    float64

	Free       float64
	Status     int
	TokenValue float64
	TokenTo    string
	InputData  string
	CreateAt   time.Time
}

func (MHCTransferLog) TableName() string {
	return "transfer_log_mhc"
}
