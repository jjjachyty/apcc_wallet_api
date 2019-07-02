package utils

import (
	"fmt"
	"math"

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
	masterChild, _ = masterChild.NewChildKey(bip32.FirstHardenedChild + 60)
	masterChild, _ = masterChild.NewChildKey(bip32.FirstHardenedChild)
	masterChild, _ = masterChild.NewChildKey(0)

	fmt.Println(math.MaxUint32)
	fmt.Println("pub:", masterChild.PublicKey())

	// masterChild, _ = masterChild.NewChildKey(100000000)
	// pubMasterChild := masterChild.PublicKey()
	// // fmt.Printf("%x  %x  %x %s\n", masterChild.Key, masterChild.ChainCode, masterChild.ChildNumber, masterChild)
	// fmt.Printf("%x\n", pubMasterChild.Key)
	// xpub := masterKey.PublicKey()

	return masterKey
}
