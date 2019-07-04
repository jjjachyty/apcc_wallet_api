package dimMod

type DimCoin struct {
	UUID     string `xorm:"varchar(36) notnull unique pk 'uuid'"`
	NameEn   string
	NameCn   string
	Symbol   string
	PriceCny float64 `json:"PriceCny,string"`
	PriceUsd float64 `json:"PriceUsd,string"`
	Icon     string
	WebSite  string
	State    string
	Synopsis string
}

func (DimCoin) TableName() string {
	return "dim_coin"
}
