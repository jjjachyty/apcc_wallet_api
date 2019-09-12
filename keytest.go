package main

import (
	"apcc_wallet_api/utils"

	"github.com/ethereum/go-ethereum/common"

	bip39 "github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"

	"fmt"

	"github.com/ethereum/go-ethereum/crypto"

	bip32 "github.com/tyler-smith/go-bip32"
)

func GetKey() {
	//拒 计 精 荷 酸 压 悲 尤 针 过 凡 纪 东 谓 嫩 哭 敢 韦 洞 穆 叛 柳 趋 瑞
	bip39.SetWordList(wordlists.ChineseSimplified)
	// 生成随机数
	// entropy, _ := bip39.NewEntropy(128)

	// 生成助记词
	// mmic, _ := bip39.NewMnemonic(entropy)
	seed := bip39.NewSeed("我 的 名 字 叫 张 力 我 爱 汪 豆 豆", "")
	fmt.Printf("seed=%x\n", seed)
	masterKey, _ := bip32.NewMasterKey(seed)
	fmt.Printf("masterkey=%x\n", masterKey.Key)
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
func main() {
	GetKey()
}
