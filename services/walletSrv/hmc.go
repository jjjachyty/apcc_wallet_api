package walletSrv

import (
	"apcc_wallet_api/models"
	"apcc_wallet_api/models/assetMod"
	"apcc_wallet_api/utils"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/crypto"
)

var hotWalletPrivateKey = "68f80940c98719851873e9e41e28f9a98f15e73df918401c51f9a000e7d84db5"
var mhcClient *ethclient.Client
var chainID *big.Int
var gasLimit = uint64(21000) // in units
var mhcAddress = "ws://119.3.108.19:8112"

func initMHCVars() {
	var err error
	if mhcAddress, err = utils.HGet("mhc", "ws_address"); err == nil {
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
	go SubscribeMHCNewHead()

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
func Logs() {
	ctx := context.Background()
	var err error
	var lastBlock *types.Block
	var block *types.Block
	var msg types.Message
	if lastBlock, err = mhcClient.BlockByNumber(ctx, nil); err == nil {
		go SubscribeMHCNewHead()
		var syncedBlock int64
		if err = models.SQLBean(&syncedBlock, "select IFNULL(max(block_number),0) from transfer_log_mhc"); err == nil {
			utils.SysLog.Debugf("开始同步区块%d至%d", syncedBlock, lastBlock.Number().Int64())
			for index := int64(syncedBlock + 1); index < lastBlock.Number().Int64(); index++ {
				mhcTransferLogs := make([]assetMod.MHCTransferLog, 0)

				if block, err = mhcClient.BlockByNumber(ctx, big.NewInt(index)); err == nil {
					if block.Transactions().Len() > 0 {
						for _, tx := range block.Transactions() {
							if msg, err = tx.AsMessage(types.NewEIP155Signer(big.NewInt(3333))); err == nil {
								floatValue, _ := big.NewFloat(0).SetString(tx.Value().String())
								value := big.NewFloat(0).Quo(floatValue, big.NewFloat(1000000000000000000))
								valueFloat, _ := value.Float64()
								// err := models.Create(
								mhcTransferLogs = append(mhcTransferLogs, assetMod.MHCTransferLog{
									TxHash:      tx.Hash().Hex(),
									BlockNumber: block.Number().Int64(),
									BlockHash:   block.Hash().Hex(),
									From:        msg.From().Hex(),
									To:          msg.To().Hex(),
									Nonce:       int64(tx.Nonce()),
									Value:       valueFloat,
									Gas:         int64(tx.Gas()),
									GasPrice:    tx.GasPrice().Int64(),
									GasUsed:     int64(block.GasUsed()),
									Status:      1,
									CreateAt:    time.Unix(int64(block.Time()), 0),
									InputData:   common.ToHex(msg.Data()),
								})
								// )
							}
						}
						//批量插入交易
						err = models.Create(mhcTransferLogs)
					}
					if err != nil {
						utils.SysLog.Errorf("同步区块%d交易至数据库出错", block.Number().Int64())
						break
					}
				}
			}

		}
	}
}

//SubscribeMHCNewHead 订阅MHC区块变化
func SubscribeMHCNewHead() {
	header := make(chan *types.Header)
	ctx := context.Background()
	if sub, err := mhcClient.SubscribeNewHead(ctx, header); err == nil {
		utils.SysLog.Debugln("开始监听区块变化")

		go func() {
			for {
				select {
				case <-sub.Err():
					return
				case hd := <-header:
					utils.SysLog.Debugf("收到新的区块%d", hd.Number.Int64())

					if block, err := mhcClient.BlockByHash(ctx, hd.Hash()); err == nil {
						mhcTransferLogs := make([]assetMod.MHCTransferLog, 0)
						if block.Transactions().Len() > 0 {
							for _, tx := range block.Transactions() {

								if msg, err := tx.AsMessage(types.NewEIP155Signer(big.NewInt(3333))); err == nil {
									floatValue, _ := big.NewFloat(0).SetString(tx.Value().String())
									value := big.NewFloat(0).Quo(floatValue, big.NewFloat(1000000000000000000))
									valueFloat, _ := value.Float64()
									fmt.Println(tx.Nonce(), block.Nonce())
									// err := models.Create(
									mhcTransferLogs = append(mhcTransferLogs, assetMod.MHCTransferLog{
										TxHash:      tx.Hash().Hex(),
										BlockNumber: block.Number().Int64(),
										BlockHash:   block.Hash().Hex(),
										From:        msg.From().Hex(),
										To:          msg.To().Hex(),
										Nonce:       int64(tx.Nonce()),
										Value:       valueFloat,
										Gas:         int64(tx.Gas()),
										GasPrice:    tx.GasPrice().Int64(),
										GasUsed:     int64(block.GasUsed()),
										Status:      1,
										CreateAt:    time.Unix(int64(block.Time()), 0),
										InputData:   common.ToHex(msg.Data()),
									})
									// )
								}
							}
							//批量插入交易
							err = models.Create(mhcTransferLogs)
						}
						if err != nil {
							utils.SysLog.Errorf("同步区块%d交易至数据库出错", hd.Number.Int64())
							break
						}
					}
				}
			}
		}()
	}
}
