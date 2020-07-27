package client

import (
	"context"
	"log"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"gopkg.in/ArcadiaMediaGroup/zcoin-rosetta-node.v0/configuration"
)

const (
	Transfer      = "transfer"
	StatusSuccess = "success"
	StatusFail    = "fail"
	ActionTypeFee = "fee"
)

const (
	WITNESS_V0 = "witness_v0_keyhash"
	P2PKH      = "pubkeyhash"
	PUBKEY     = "pubkey"
)

const BASE_CURRENCY_DECIMAL_DIVIDER = 100000000
const BASE_CURRENCY_DECIMAL_COUNT = 8
const CURRENCY_SYMBOL = "XZC"

func IsValidPaymentType(paymentType string) bool {
	return paymentType == P2PKH || paymentType == WITNESS_V0 || paymentType == PUBKEY
}

// ZcoinClientRPC is an implementation of ZcoinClient using RPC.
type ZcoinClientRPC struct {
	rpcConnConfig     *rpcclient.ConnConfig
	applicationConfig *configuration.Config
}

// NewZcoinClient returns an implementation of ZcoinClient
func NewZcoinClient(applicationConfig *configuration.Config) (cli ZcoinClient) {
	rpcConnConfig := rpcclient.ConnConfig{
		Host:         applicationConfig.Node.Endpoint,
		User:         applicationConfig.Node.Username,
		Pass:         applicationConfig.Node.Password,
		HTTPPostMode: true,                               // Bitcoin core only supports HTTP POST mode
		DisableTLS:   !applicationConfig.Node.TLSEnabled, // Bitcoin core does not provide TLS by default
	}

	return &ZcoinClientRPC{
		rpcConnConfig:     &rpcConnConfig,
		applicationConfig: applicationConfig,
	}
}

func (rpcClient *ZcoinClientRPC) reconnect() (client *rpcclient.Client) {
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(rpcClient.rpcConnConfig, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

// GetConfig retrieves the general application config that has been configured
func (rpcClient *ZcoinClientRPC) GetConfig() *configuration.Config {
	return rpcClient.applicationConfig
}

// GetStatus will return the Blockchain base info based on that node
func (rpcClient *ZcoinClientRPC) GetStatus(ctx context.Context) (*btcjson.GetBlockChainInfoResult, error) {
	client := rpcClient.reconnect()
	defer client.Shutdown()

	result, err := client.GetBlockChainInfo()
	return result, err
}

// GetBlock will return you the block specification for a given height
func (rpcClient *ZcoinClientRPC) GetBlock(ctx context.Context, height int64) (*btcjson.GetBlockVerboseResult, error) {
	client := rpcClient.reconnect()
	defer client.Shutdown()

	result, err := client.GetBlockHash(height)
	block, err := client.GetBlockVerbose(result)
	return block, err
}

// GetBlockByHash will return you the block specification for a given height
func (rpcClient *ZcoinClientRPC) GetBlockByHash(ctx context.Context, hash string) (*btcjson.GetBlockVerboseResult, error) {
	client := rpcClient.reconnect()
	defer client.Shutdown()

	blockHash, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		return nil, err
	}

	block, err := client.GetBlockVerbose(blockHash)
	if err != nil {
		return nil, err
	}

	return block, nil
}

// GetLatestBlock returns the latest Zcoin block.
func (rpcClient *ZcoinClientRPC) GetLatestBlock(ctx context.Context) (*wire.MsgBlock, error) {
	client := rpcClient.reconnect()
	defer client.Shutdown()

	latestBlockHash, err := client.GetBestBlockHash()
	if err != nil {
		return nil, err
	}

	block, err := client.GetBlock(latestBlockHash)
	if err != nil {
		return nil, err
	}

	return block, nil
}

// GetBlockByHashWithTransaction returns the Zcoin block including transactions
func (rpcClient *ZcoinClientRPC) GetBlockByHashWithTransaction(ctx context.Context, hash string) (*btcjson.GetBlockVerboseTxResult, error) {
	client := rpcClient.reconnect()
	defer client.Shutdown()

	blockHash, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		return nil, err
	}

	block, err := client.GetBlockVerboseTx(blockHash)
	if err != nil {
		return nil, err
	}

	return block, nil
}
