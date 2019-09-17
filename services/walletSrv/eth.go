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
var ethHotWalletAddress string
var ethColdWalletAddress string

var ethClient *ethclient.Client

var privateKey string
var ethContractsAddress string

var instance *usdt.TetherToken
var rpcAddress string

func initETHVars() {
	var err error
	if rpcAddress, err = utils.HGet("eth", "ws_address"); err == nil {
		if ethHotWalletAddress, err = utils.HGet("eth", "hot_wallet_address"); err == nil {
			if privateKey, err = utils.HGet("eth", "hot_wallet_privatekey"); err == nil {
				if ethContractsAddress, err = utils.HGet("eth", "contracts_address_usdt"); err == nil {
					ethColdWalletAddress, err = utils.HGet("eth", "cold_wallet_address")
				}
			}
		}
	}

	if err != nil {
		utils.SysLog.Panicf("初始化eth配置出错%v", err)
	}
}

func InitETHClient() {

	var err error
	initETHVars()
	if ethClient, err = ethclient.Dial(rpcAddress); err == nil {

		if auth, err = getAuth(privateKey); err == nil {
			if instance, err = getInstance(ethContractsAddress); err != nil {
				utils.SysLog.Panicf("获取USDT实例失败 %v", err)
			}
			if balance, err := GetUSDTBalance(auth.From.Hex()); err == nil {
				utils.SysLog.Debugf("当前账户剩余%sUSDT", balance.String())
			}
		} else {
			utils.SysLog.Panic("获取AUth 失败", err)
		}

	} else {
		utils.SysLog.Panic("USDT 服务器连接失败")
	}

}

func SendUSDT(toAddress string, amount *big.Int) (*types.Transaction, error) {
	return instance.Transfer(&bind.TransactOpts{
		Signer: auth.Signer,
		From:   auth.From,
	}, common.HexToAddress(toAddress), amount)
}

func SendUSDTByPrivateKey(privateKey string, toAddress common.Address, amount *big.Int) (tx *types.Transaction, custAuth *bind.TransactOpts, err error) {

	if custAuth, err = getAuth(privateKey); err == nil {
		fmt.Println("privateKey", privateKey, "toAddress", toAddress.Hex(), "from", custAuth.From.Hex(), "amount", amount)
		tx, err = instance.Transfer(&bind.TransactOpts{
			Signer: custAuth.Signer,
			From:   custAuth.From,
		}, toAddress, amount)
		return
	}
	return nil, nil, err
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

//GetHotWalletUSDTBalance 获取热钱包USDT余额
func GetHotWalletUSDTBalance() (*big.Int, error) {
	return instance.BalanceOf(nil, common.HexToAddress(ethHotWalletAddress))
}

func getAuth(pkHex string) (*bind.TransactOpts, error) {
	var err error
	var ecdsaPrivateKey *ecdsa.PrivateKey
	var nonce uint64
	var gasPrice *big.Int
	if ecdsaPrivateKey, err = crypto.HexToECDSA(pkHex); err == nil {

		if ecdsaPublicKey, ok := ecdsaPrivateKey.Public().(*ecdsa.PublicKey); ok {
			fromAddress := crypto.PubkeyToAddress(*ecdsaPublicKey)
			if nonce, err = ethClient.PendingNonceAt(context.Background(), fromAddress); err == nil {
				if gasPrice, err = ethClient.SuggestGasPrice(context.Background()); err == nil {
					auth := bind.NewKeyedTransactor(ecdsaPrivateKey)
					auth.Nonce = big.NewInt(int64(nonce))
					auth.Value = big.NewInt(0) // in wei
					auth.GasLimit = gasLimit   // in units
					auth.GasPrice = gasPrice
					return auth, nil
				}
			}

		}
	}

	return nil, err

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
func SendETH(toAddress string, amount *big.Int) error {
	var err error
	var ecdsaPrivateKey *ecdsa.PrivateKey
	var chainID *big.Int
	var nonce uint64
	var transaction *types.Transaction
	if ecdsaPrivateKey, err = crypto.HexToECDSA(privateKey); err == nil {
		if nonce, err = ethClient.PendingNonceAt(context.Background(), auth.From); err == nil {
			transaction = types.NewTransaction(nonce, common.HexToAddress(toAddress), amount, auth.GasLimit, auth.GasPrice, nil)
			if chainID, err = ethClient.NetworkID(context.Background()); err == nil {
				transaction, err = types.SignTx(transaction, types.NewEIP155Signer(chainID), ecdsaPrivateKey)
				err = ethClient.SendTransaction(context.Background(), transaction)
			}

		}

	}
	return err
}

//GetETHBalance 获取ETh余额
func GetETHBalance(addressHex string) (*big.Int, error) {
	return ethClient.BalanceAt(context.Background(), common.HexToAddress(addressHex), nil)
}

//GetETHHotWalletAddressBalance 获取热钱包eth余额
func GetETHHotWalletAddressBalance() (*big.Int, error) {
	return ethClient.BalanceAt(context.Background(), common.HexToAddress(ethHotWalletAddress), nil)
}

//GetHotWalletAddress 获取热钱包地址
func GetETHHotWalletAddress() common.Address {
	return common.HexToAddress(ethHotWalletAddress)
}

//GetClodWalletAddress 获取USDT冷钱包地址
func GetETHClodWalletAddress() common.Address {
	return common.HexToAddress(ethColdWalletAddress)
}
