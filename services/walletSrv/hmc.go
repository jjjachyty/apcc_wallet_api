package walletSrv

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/crypto"
)

var privateKeyHex = "68f80940c98719851873e9e41e28f9a98f15e73df918401c51f9a000e7d84db5"
var mhcClient *ethclient.Client
var chainID *big.Int
var gasLimit = uint64(21000) // in units

func SendMHC(amount *big.Int, toAddressHex string) (string, string, error) {
	var err error
	var privateKey *ecdsa.PrivateKey
	var fromAddress common.Address
	var nonce uint64
	var gasPrice *big.Int
	var signedTx *types.Transaction
	toAddress := common.HexToAddress(toAddressHex)
	if privateKey, err = crypto.HexToECDSA(privateKeyHex); err == nil {
		if fromAddress, err = getETHAddressByPK(privateKey); err == nil {
			if nonce, err = mhcClient.PendingNonceAt(context.Background(), fromAddress); err == nil {
				if gasPrice, err = mhcClient.SuggestGasPrice(context.Background()); err == nil {
					tx := types.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, nil)
					if signedTx, err = types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey); err == nil {
						if err = mhcClient.SendTransaction(context.Background(), signedTx); err == nil {
							return fromAddress.Hex(), signedTx.Hash().Hex(), nil
						}

					}
				}

			}
		}
	}
	return "", "", err
}

func GetTxsByHashHex(txsHex string) (*types.Transaction, bool, error) {
	return mhcClient.TransactionByHash(context.Background(), common.HexToHash(txsHex))
}

func init() {
	client, err := ethclient.Dial("http://119.3.108.19:8110")
	if err != nil {
		panic("http://119.3.108.19:8110 MHC客户端创建失败")
	}
	mhcClient = client

	id, err := client.NetworkID(context.Background())
	if err != nil {
		panic("获取ChainID 失败")
	}
	chainID = id
}

func getETHAddressByPK(privateKey *ecdsa.PrivateKey) (common.Address, error) {
	var err error
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		err = errors.New("获取地址失败")
	}
	return crypto.PubkeyToAddress(*publicKeyECDSA), err
}
