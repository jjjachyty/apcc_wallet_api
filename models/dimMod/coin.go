package dimMod

type DimCoin struct {
	UUID             string `xorm:"varchar(36) notnull unique pk 'uuid'"`
	NameEn           string
	NameCn           string
	Symbol           string
	PriceCny         float64 `json:"PriceCny,string"`
	PriceUsd         float64 `json:"PriceUsd,string"`
	ExchangeUsdtFree float64 //USDT 转账费用
	TransferFree     float64 //同币种转账费用
	Icon             string
	WebSite          string
	State            int
	Synopsis         string
}

func (DimCoin) TableName() string {
	return "dim_coin"
}
