package walletSrv

import (
	"apcc_wallet_api/utils"
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/crypto"
)

var hotWalletPrivateKey = "68f80940c98719851873e9e41e28f9a98f15e73df918401c51f9a000e7d84db5"
var mhcClient *ethclient.Client
var chainID *big.Int
var gasLimit = uint64(21000) // in units
var mhcAddress = "http://119.3.108.19:8110"

func initMHCVars() {
	var err error
	if mhcAddress, err = utils.HGet("mhc", "rpc_address"); err == nil {
		hotWalletPrivateKey, err = utils.HGet("mhc", "hot_wallet_privatekey")
	}

	if err != nil {
		utils.SysLog.Panicf("初始化mhc配置出错%v", err)
	}
}

func InitMHCClient() {
	initMHCVars()
	client, err := ethclient.Dial(mhcAddress)
	if err != nil {
		utils.SysLog.Panicf("MHC客户端%s创建失败", mhcAddress)
	}
	mhcClient = client

	id, err := client.NetworkID(context.Background())
	if err != nil {
		utils.SysLog.Panic("获取ChainID 失败")
	}
	chainID = id
	utils.SysLog.Debugln("MHC客户端初始化成功")

}

func SendMHC(amount *big.Int, toAddressHex string) (address string, txs string, err error) {
	var privateKey *ecdsa.PrivateKey
	var fromAddress common.Address
	var nonce uint64
	var gasPrice *big.Int
	var signedTx *types.Transaction
	toAddress := common.HexToAddress(toAddressHex)
	if privateKey, err = crypto.HexToECDSA(hotWalletPrivateKey); err == nil {
		if fromAddress, err = GetETHAddressByPK(privateKey); err == nil {
			if nonce, err = mhcClient.PendingNonceAt(context.Background(), fromAddress); err == nil {
				if gasPrice, err = mhcClient.SuggestGasPrice(context.Background()); err == nil {
					tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, nil)
					if signedTx, err = types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey); err == nil {
						if err = mhcClient.SendTransaction(context.Background(), signedTx); err == nil {
							address, txs, err = fromAddress.Hex(), signedTx.Hash().Hex(), nil
							return address, txs, err
						}

					}
				}

			}
		}
	}
	return "", "", err
}

func SendMHCByPrivateKey(privateKeyHex string, amount *big.Int, toAddressHex string) (address string, tx *types.Transaction, err error) {
	var privateKey *ecdsa.PrivateKey
	var fromAddress common.Address
	var nonce uint64
	var gasPrice *big.Int
	// var tx *types.Transaction
	toAddress := common.HexToAddress(toAddressHex)
	if privateKey, err = crypto.HexToECDSA(privateKeyHex); err == nil {
		if fromAddress, err = GetETHAddressByPK(privateKey); err == nil {
			if nonce, err = mhcClient.PendingNonceAt(context.Background(), fromAddress); err == nil {
				if gasPrice, err = mhcClient.SuggestGasPrice(context.Background()); err == nil {
					tx = types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, nil)
					if tx, err = types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey); err == nil {
						if err = mhcClient.SendTransaction(context.Background(), tx); err == nil {

							address, err = fromAddress.Hex(), nil
							return
						}

					}
				}

			}
		}
	}
	return "", nil, err
}

func GetTxsByHashHex(txsHex string) (*types.Transaction, bool, error) {
	return mhcClient.TransactionByHash(context.Background(), common.HexToHash(txsHex))
}

func GetETHAddressByPKHex(privateKeyHex string) (address common.Address, err error) {
	var privateKey *ecdsa.PrivateKey

	if privateKey, err = crypto.HexToECDSA(privateKeyHex); err == nil {
		return GetETHAddressByPK(privateKey)
	}
	return
}
func GetETHAddressByPK(privateKey *ecdsa.PrivateKey) (address common.Address, err error) {
	publicKey := privateKey.Public()
	if publicKey, ok := publicKey.(*ecdsa.PublicKey); ok {
		return crypto.PubkeyToAddress(*publicKey), nil
	}

	return
}

func GetMHCBalance(address string) (*big.Int, error) {
	return mhcClient.BalanceAt(context.Background(), common.HexToAddress(address), nil)
}
func GetMHCTransactionReceipt(txHash string) (*types.Receipt, error) {

	return mhcClient.TransactionReceipt(context.Background(), common.HexToHash(txHash))

}
func GetMHCLastBlockNum() int64 {
	if header, err := mhcClient.HeaderByNumber(context.Background(), nil); err == nil {
		return header.Number.Int64()
	}
	return 1
}
func GetTxs(address string) (*big.Int, error) {
	return mhcClient.BalanceAt(context.Background(), common.HexToAddress(address), nil)
}

func GetMHCGas() *big.Int {
	if gasPrice, err := mhcClient.SuggestGasPrice(context.Background()); err == nil {
		return new(big.Int).Mul(gasPrice, big.NewInt(2100))
	}
	return new(big.Int).Mul(big.NewInt(100000000000), big.NewInt(2100))
}
