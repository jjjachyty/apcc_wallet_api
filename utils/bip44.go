package utils

import (
	"crypto/sha256"
	"fmt"
	"math"

	"github.com/btcsuite/btcutil/base58"

	"golang.org/x/crypto/ripemd160"

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

//将公钥进行两次哈希
func Ripemd160Hash(publicKey []byte) []byte {

	//1.hash256
	hash256 := sha256.New()
	hash256.Write(publicKey)
	hash := hash256.Sum(nil)

	//2.ripemd160
	ripemd160 := ripemd160.New()
	ripemd160.Write(hash)

	return ripemd160.Sum(nil)
}

//两次sha256哈希生成校验和
func CheckSum(bytes []byte) []byte {

	//hasher := sha256.New()
	//hasher.Write(bytes)
	//hash := hasher.Sum(nil)
	//与下面一句等同
	//hash := sha256.Sum256(bytes)

	hash1 := sha256.Sum256(bytes)
	hash2 := sha256.Sum256(hash1[:])

	return hash2[:4]
}

func GetAddress(pubKeyByts []byte) string {

	//1.使用RIPEMD160(SHA256(PubKey)) 哈希算法，取公钥并对其哈希两次
	ripemd160Hash := Ripemd160Hash(pubKeyByts)

	//2.拼接版本
	version_ripemd160Hash := append([]byte{0x00}, ripemd160Hash...)
	//3.两次sha256生成校验和
	checkSumBytes := CheckSum(version_ripemd160Hash)
	//4.拼接校验和
	bytes := append(version_ripemd160Hash, checkSumBytes...)

	return base58.Encode(bytes)
}
