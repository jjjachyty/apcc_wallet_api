package walletSrv

import (
	"apcc_wallet_api/models/assetMod"
	"apcc_wallet_api/services/assetSrv"
	"apcc_wallet_api/utils"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	bip39 "github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"

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
	btcAccountPrivateKey = "xprvA1wLVjqnGWur8JR1Amx5tUSdAdWbMFGjWg1x6QxyhKyve1vY7Y8eyNeo4x3ZDE7Jos59QGuHhUZ3Cc1X1mHyx5TxCunS5g8KoibLiXidQcp"
	ethAccountPrivateKey = "xprvA1wLVjqnGWur8JR1Amx5tUSdAdWbMFGjWg1x6QxyhKyve1vY7Y8eyNeo4x3ZDE7Jos59QGuHhUZ3Cc1X1mHyx5TxCunS5g8KoibLiXidQcp"
	btcAccountPub        = "xpub6EELkV1BFXMvWtH7XsbGTnojnd1CEFgJecSjo5qs2s7Q88yDjBtn9Rc6xTTchWpYYNtYBPJFNny1UzHQCCEX9X2ZzdxjqZCMWWo69MNfQ8P"
	ethAccountPub        = "xpub6EvguFNg6tU9LnVUGoV6FcPMifM5khzastwYtoNbFfWuWpFgf5SuXAyGvEzuoa7qXsnW1UmtBFoUVJCBky9EwcRHaohd1ZEXNvsNMnmBdiH"
)

//获取BTC地址
func GetETHAddressPrivateKey(acountID uint32) ([]byte, error) {
	key, err := bip32.B58Deserialize(ethAccountPrivateKey)
	if err == nil {
		if child, err := key.NewChildKey(acountID); err == nil {
			return child.Key, nil
		}
	}
	return nil, err
}

//获取BTC地址
func GetBtcAddress(acountID uint32) (string, error) {
	key, err := bip32.B58Deserialize(btcAccountPub)
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
	key, err := bip32.B58Deserialize(ethAccountPub)
	if err == nil {
		if child, err := key.NewChildKey(acountID); err == nil {

			pubKey := utils.ExpandPublicKey(child.Key)
			return crypto.PubkeyToAddress(*pubKey).Hex(), err
		}
	}
	return "", err
}

func GetAddress(userid string, acountID uint32) ([]assetMod.Asset, error) {
	var assets = make([]assetMod.Asset, 1)
	var ethAddr string
	var err error

	if ethAddr, err = GetEthAddress(acountID); err == nil {
		assets[0] = assetMod.Asset{UUID: userid, Symbol: assetSrv.USDT_COIN_SYMBOL, BaseOn: BASE_ON_ETH, Address: ethAddr}
		// assets[1] = assetMod.Asset{UUID: userid, Symbol: assetSrv.USDT_COIN_SYMBOL, BaseOn: BASE_ON_ETH, Address: ethAddr}
	}
	// if btcAddr, err = GetBtcAddress(acountID); err == nil {

	// 	assets[2] = assetMod.Asset{UUID: userid, Symbol: assetSrv.USDT_COIN_SYMBOL, BaseOn: BASE_ON_BTC, Address: btcAddr}

	// }
	return assets, err
}

func GetKey() {
	//拒 计 精 荷 酸 压 悲 尤 针 过 凡 纪 东 谓 嫩 哭 敢 韦 洞 穆 叛 柳 趋 瑞
	bip39.SetWordList(wordlists.ChineseSimplified)
	// 生成随机数
	// entropy, _ := bip39.NewEntropy(128)

	// 生成助记词
	// mmic, _ := bip39.NewMnemonic(entropy)
	seed := bip39.NewSeed("拒 计 精 荷 酸 压 悲 尤 针 过 凡 纪 东 谓 嫩 哭 敢 韦 洞 穆 叛 柳 趋 瑞", "")

	masterKey, _ := bip32.NewMasterKey(seed)
	purposeKey, _ := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	coinTypeKey, _ := purposeKey.NewChildKey(bip32.FirstHardenedChild + 3333)
	accountKey, _ := coinTypeKey.NewChildKey(bip32.FirstHardenedChild + 0)
	changeKey, _ := accountKey.NewChildKey(0)
	println("accountPK=", changeKey.B58Serialize(), "accountPubK=", changeKey.PublicKey().B58Serialize())
	addressKey, _ := changeKey.NewChildKey(1)
	fmt.Println("address pk= ", common.ToHex(addressKey.Key))
	fmt.Printf("PublicKey().Key %x    %x", addressKey.PublicKey().Key, addressKey.Key)
	pubKey := utils.ExpandPublicKey(addressKey.PublicKey().Key)
	fmt.Println("第一个公钥地址", crypto.PubkeyToAddress(*pubKey).Hex())
	fmt.Println("第一个私钥地址", crypto.PubkeyToAddress(*pubKey).Hex())

}
