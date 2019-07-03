package utils

import (
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"

	btcutil "github.com/FactomProject/btcutilecc"
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

// As described at https://crypto.stackexchange.com/a/8916
func ExpandPublicKey(key []byte) *ecdsa.PublicKey {
	Y := big.NewInt(0)
	X := big.NewInt(0)
	X.SetBytes(key[1:])
	cure := btcutil.Secp256k1()

	curveParams := cure.Params()
	// y^2 = x^3 + ax^2 + b
	// a = 0
	// => y^2 = x^3 + b
	ySquared := big.NewInt(0)
	ySquared.Exp(X, big.NewInt(3), nil)
	ySquared.Add(ySquared, curveParams.B)

	Y.ModSqrt(ySquared, curveParams.P)

	Ymod2 := big.NewInt(0)
	Ymod2.Mod(Y, big.NewInt(2))

	signY := uint64(key[0]) - 2
	if signY != Ymod2.Uint64() {
		Y.Sub(curveParams.P, Y)
	}

	return &ecdsa.PublicKey{cure, X, Y}
}
