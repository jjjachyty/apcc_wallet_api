package dimMod

type DimCoin struct {
	UUID     string `xorm:"varchar(36) notnull unique pk 'uuid'"`
	NameEn   string
	NameCn   string
	Symbol   string
	PriceCny float64
	PriceUsd float64
	Icon     string
	WebSite  string
	State    string
	Synopsis string
}

func (DimCoin) TableName() string {
	return "dim_coin"
}
