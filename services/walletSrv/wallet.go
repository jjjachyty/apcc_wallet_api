package walletSrv

import (
	"apcc_wallet_api/models/assetMod"
	"apcc_wallet_api/services/assetSrv"
	"apcc_wallet_api/utils"
	"fmt"

	"github.com/btcsuite/btcd/wire"

	"github.com/btcsuite/btcd/chaincfg"

	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/crypto"

	bip32 "github.com/tyler-smith/go-bip32"
)

const (
	BASE_ON_BTC = "BTC"
	BASE_ON_ETH = "ETH"
)

const (
	// mnemonic  = "拒 计 精 荷 酸 压 悲 尤 针 过 凡 纪 东 谓 嫩 哭 敢 韦 洞 穆 叛 柳 趋 瑞"
	// passwd    = "123456"
	// masterkey = "xprv9s21ZrQH143K2BHJXFAdkZhaUsc1qfuTnd1RdBQFLm8oxx5uD9NAQem3H9QpKmqFMvHGr1ypjSHFmPBWHcev7XSydBbBu214NzNcywvDTpJ"
	accountBtcPub = "xpub6EELkV1BFXMvWtH7XsbGTnojnd1CEFgJecSjo5qs2s7Q88yDjBtn9Rc6xTTchWpYYNtYBPJFNny1UzHQCCEX9X2ZzdxjqZCMWWo69MNfQ8P"
	accountEthPub = "xpub6FJMotUfHzsuxkB7jLZ65XMjDcWfpYcB2LgGmT2buLBorDXcknGdVzWAcyMmYu49Kv7sEcUbh4XW83giKx6e6J2WfhXhCBF5PZowEvARVQR"
)

//获取BTC地址
func GetBtcAddress(acountID uint32) (string, error) {
	key, err := bip32.B58Deserialize(accountBtcPub)
	if err == nil {
		if child, err := key.NewChildKey(acountID); err == nil {

			ext, _ := hdkeychain.NewKeyFromString(child.String())
			address, _ := ext.Address(&chaincfg.Params{Net: wire.MainNet})
			return address.EncodeAddress(), err
		}
	}
	return "", err
}

//获取ETH地址
func GetEthAddress(acountID uint32) (string, error) {
	key, err := bip32.B58Deserialize(accountEthPub)
	if err == nil {
		if child, err := key.NewChildKey(acountID); err == nil {

			pubKey := utils.ExpandPublicKey(child.Key)
			fmt.Printf("%x", child.ChildNumber)
			return crypto.PubkeyToAddress(*pubKey).Hex(), err
		}
	}
	return "", err
}

func GetAddress(userid string, acountID uint32) ([]assetMod.Asset, error) {
	var assets = make([]assetMod.Asset, 3)
	var ethAddr string
	var err error

	if ethAddr, err = GetEthAddress(acountID); err == nil {
		assets[0] = assetMod.Asset{UUID: userid, Symbol: assetSrv.HMC_COIN_SYMBOL, BaseOn: BASE_ON_ETH, Address: ethAddr}
		assets[1] = assetMod.Asset{UUID: userid, Symbol: assetSrv.USDT_COIN_SYMBOL, BaseOn: BASE_ON_ETH, Address: ethAddr}
	}
	// if btcAddr, err = GetBtcAddress(acountID); err == nil {

	// 	assets[2] = assetMod.Asset{UUID: userid, Symbol: assetSrv.USDT_COIN_SYMBOL, BaseOn: BASE_ON_BTC, Address: btcAddr}

	// }
	return assets, err
}
