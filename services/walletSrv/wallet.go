package walletSrv

import (
	"github.com/tyler-smith/go-bip32"
	"fmt"

	hdwallet "github.com/wemeetagain/go-hdwallet"
)

const (
	mnemonic  = "拒 计 精 荷 酸 压 悲 尤 针 过 凡 纪 东 谓 嫩 哭 敢 韦 洞 穆 叛 柳 趋 瑞"
	passwd    = "123456"
	masterkey = "xprv9s21ZrQH143K2BHJXFAdkZhaUsc1qfuTnd1RdBQFLm8oxx5uD9NAQem3H9QpKmqFMvHGr1ypjSHFmPBWHcev7XSydBbBu214NzNcywvDTpJ"
	masterpub = "xpub661MyMwAqRbcEfMmdGhe7heK2uSWF8dK9qw2RZoru6fnqkR3kggQxT5X8RaF9eJyn1uNVZSmXo8BqYfCn1syaZq5UsVAMRuTVZGsttfeHjz"
)

func GetAddress(acountID string) {

}

func getMasterPub() {
   
	
	// Generate a random 256 bit seed
	seed, err := hdwallet.GenSeed(256)

	// Create a master private key
	masterprv := hdwallet.MasterKey(seed)

	// Convert a private key to public key
	masterpub := masterprv.Pub()
	masterpub.
	// Generate new child key based on private or public key
	childprv, err := masterprv.Child(0)
	childpub, err := masterpub.Child(0)
	fmt.Println(childprv, childpub)
	// Create bitcoin address from public key
	address := childpub.Address()
	fmt.Println(address)
	// Convenience string -> string Child and ToAddress functions
	walletstring := childpub.String()
	childstring, err := hdwallet.StringChild(walletstring, 0)
	childaddress, err := hdwallet.StringAddress(childstring)
	fmt.Println(childaddress, err)
}
