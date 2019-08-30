package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var client *ethclient.Client
var err error

func main() {

	client, err = ethclient.Dial("http://119.3.108.19:8110")
	send()
	Bal()
}
func send() {
	privateKey, err := crypto.HexToECDSA("68f80940c98719851873e9e41e28f9a98f15e73df918401c51f9a000e7d84db5")
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	value, _ := big.NewInt(0).SetString("1000000000000000000000", 0) // in wei (1 eth)
	gasLimit := uint64(21000)                                        // in units
	// gasPrice := big.NewInt(30000000000)        // in wei (30 gwei)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	toAddress := common.HexToAddress("0x88761000d7fb6080490d54800fe5252e1a35d84d")
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex()) // tx sent: 0x77006fcb3938f648e2cc65bafd27dec30b9bfbe9df41f78498b9c8b7322a249e

}
func Bal() {
	bal, _ := client.BalanceAt(context.Background(), common.HexToAddress("0x52D0D74F73D1d4D5D054AB588e17E2EC9e5CDcbB"), nil)
	fmt.Println(bal.String())
}
