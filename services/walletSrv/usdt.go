package walletSrv

import (
	"apcc_wallet_api/services/walletSrv/usdt"
	"apcc_wallet_api/utils"
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var auth *bind.TransactOpts
var client *ethclient.Client
var privateKey = "68f80940c98719851873e9e41e28f9a98f15e73df918401c51f9a000e7d84db5"
var usdtContractsAddress = "0x27afD2BC77Ec8A520aBaf6BE2257A251733e24b3"
var instance *usdt.TetherToken

func init() {
	var err error
	client, err = ethclient.Dial("http://119.3.108.19:8110")
	if err == nil {
		auth = getAuth(privateKey)
		if instance, err = getInstance(usdtContractsAddress); err != nil {
			utils.SysLog.Panicf("获取USDT实例失败 %v", err)
		}

	} else {
		panic("USDT 服务器连接失败")
	}
}

func SendUSDT(toAddress string, amount *big.Int) (*types.Transaction, error) {
	return instance.Transfer(nil, common.HexToAddress(toAddress), amount)
}

func getAuth(pkHex string) *bind.TransactOpts {
	privateKey, err := crypto.HexToECDSA(pkHex)
	if err != nil {
		utils.SysLog.Panic(err)
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

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice
	return auth
}

func getInstance(contractAddress string) (*usdt.TetherToken, error) {
	address := common.HexToAddress(contractAddress)
	return usdt.NewTetherToken(address, client)

}
