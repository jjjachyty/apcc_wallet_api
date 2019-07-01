package utils

import (
	"fmt"

	bip32 "github.com/tyler-smith/go-bip32"
	bip39 "github.com/tyler-smith/go-bip39"
	"github.com/tyler-smith/go-bip39/wordlists"
)

func GetMnemonic(bitSize int) string {
	bip39.SetWordList(wordlists.ChineseSimplified)

	entropy, _ := bip39.NewEntropy(bitSize)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	return mnemonic
}

func GetMasterKey(mnemonic string, passwd string) *bip32.Key {
	// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	seed := bip39.NewSeed(mnemonic, passwd)

	masterKey, _ := bip32.NewMasterKey(seed)
	masterChild, _ := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	masterChild, _ = masterChild.NewChildKey(bip32.FirstHardenedChild)
	masterChild, _ = masterChild.NewChildKey(bip32.FirstHardenedChild)
	masterChild, _ = masterChild.NewChildKey(0)

	masterChild, _ = masterChild.NewChildKey(0)
	pubMasterChild := masterChild.PublicKey()
	// fmt.Printf("%x  %x  %x %s\n", masterChild.Key, masterChild.ChainCode, masterChild.ChildNumber, masterChild)
	fmt.Printf("%x\n", pubMasterChild.Key)
	// xpub := masterKey.PublicKey()

	key, _ := bip32.B58Deserialize("xpub661MyMwAqRbcEfMmdGhe7heK2uSWF8dK9qw2RZoru6fnqkR3kggQxT5X8RaF9eJyn1uNVZSmXo8BqYfCn1syaZq5UsVAMRuTVZGsttfeHjz")

	childKey, err := key.NewChildKey(44)
	// fmt.Printf("%x  %x  %x %s\n", childKey.Key, childKey.ChainCode, childKey.ChildNumber, childKey)

	fmt.Println(err)
	childKey, _ = childKey.NewChildKey(0)
	childKey, _ = childKey.NewChildKey(0)
	// fmt.Printf("%x  %x  %x %s\n", childKey.Key, childKey.ChainCode, childKey.ChildNumber, childKey)

	childKey, _ = childKey.NewChildKey(0)
	// fmt.Printf("%x  %x  %x %s\n", childKey.Key, childKey.ChainCode, childKey.ChildNumber, childKey)
	childKey, _ = childKey.NewChildKey(0)
	// fmt.Printf("%x  %x  %x %s", childKey.Key, childKey.ChainCode, childKey.ChildNumber, childKey)

	fmt.Printf("%x\n", childKey.PublicKey().Key)
	return masterKey
}
