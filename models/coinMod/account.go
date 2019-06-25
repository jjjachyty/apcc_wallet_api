package coinMod

import "math/big"

type Account struct {
	Owner           string //所属人
	CoinSymbol      string
	CoinName        string
	CoinCode        string
	Address         string //地址
	CurrentBalance  big.Int
	FreezingBalance big.Int
}
