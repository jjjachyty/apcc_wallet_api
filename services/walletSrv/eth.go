package walletSrv

import (
	"apcc_wallet_api/services/walletSrv/usdt"
	"apcc_wallet_api/utils"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var auth *bind.TransactOpts
var ethClient *ethclient.Client
var privateKey = "68f80940c98719851873e9e41e28f9a98f15e73df918401c51f9a000e7d84db5"
var usdtContractsAddress = "0x27afD2BC77Ec8A520aBaf6BE2257A251733e24b3"
var instance *usdt.TetherToken
var usdtAddress = "http://119.3.108.19:8110"

func InitUSDTClient() {
	var err error
	ethClient, err = ethclient.Dial(usdtAddress)

	if err == nil {
		auth = getAuth(privateKey)
		if instance, err = getInstance(usdtContractsAddress); err != nil {
			utils.SysLog.Panicf("获取USDT实例失败 %v", err)
		}
		if balance, err := GetUSDTBalance(auth.From.Hex()); err == nil {
			utils.SysLog.Debugf("当前账户剩余%sUSDT", balance.String())
		}

	} else {
		utils.SysLog.Panic("USDT 服务器连接失败")
	}

	utils.SysLog.Debugln("USDT客户端初始化成功")
}

func SendUSDT(toAddress string, amount *big.Int) (*types.Transaction, error) {
	return instance.Transfer(&bind.TransactOpts{
		Signer: auth.Signer,
		From:   auth.From,
	}, common.HexToAddress(toAddress), amount)
}

func SendUSDTByPrivateKey(privateKey string, toAddress string, amount *big.Int) (*types.Transaction, error) {
	auth := getAuth(privateKey)
	fmt.Println("privateKey", privateKey, "toAddress", toAddress, "from", auth.From.Hex(), "amount", amount)
	return instance.Transfer(&bind.TransactOpts{
		Signer: auth.Signer,
		From:   auth.From,
	}, common.HexToAddress(toAddress), amount)
}

func GetAuth() *bind.TransactOpts {
	return auth
}

func GetGas() *big.Int {
	if gas, err := ethClient.SuggestGasPrice(context.Background()); err == nil {
		return new(big.Int).Mul(gas, big.NewInt(int64(gasLimit)))
	}
	return big.NewInt(21000000000000)
}
func GetUSDTBalance(address string) (*big.Int, error) {
	return instance.BalanceOf(nil, common.HexToAddress(address))
}

func getAuth(pkHex string) *bind.TransactOpts {
	privateKey, err := crypto.HexToECDSA(pkHex)
	if err != nil {
		utils.SysLog.Errorln("私钥转换ecdsa.PrivateKey出错")
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		utils.SysLog.Errorln("公钥转换*ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := ethClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		utils.SysLog.Errorln("获取当前地址Nonce出错")
	}

	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	fmt.Println("gasPrice=================================", gasPrice.String())
	if err != nil {
		utils.SysLog.Errorln("获取gasPrice出错")

	}
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimit   // in units
	auth.GasPrice = gasPrice
	return auth
}

func getInstance(contractAddress string) (*usdt.TetherToken, error) {
	address := common.HexToAddress(contractAddress)
	return usdt.NewTetherToken(address, ethClient)

}
func GetETHTransaction(txHash string) (*types.Transaction, bool, error) {

	return ethClient.TransactionByHash(context.Background(), common.HexToHash(txHash))

}
func GetETHTransactionReceipt(txHash string) (*types.Receipt, error) {

	return ethClient.TransactionReceipt(context.Background(), common.HexToHash(txHash))

}
func GetETHLastBlockNum() int64 {

	if header, err := ethClient.HeaderByNumber(context.Background(), nil); err == nil {
		return header.Number.Int64()
	}
	return 1
}
func SendETH(toAddress string) int64 {
	if nonce, err := ethClient.PendingNonceAt(context.Background(), auth.From); err == nil {
		tx := types.NewTransaction(nonce, common.HexToAddress(toAddress), value, auth.GasLimit, auth.GasPrice, nil)
	}

}
