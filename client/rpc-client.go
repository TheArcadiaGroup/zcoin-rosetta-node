package client

import (
	"context"
	"log"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/marpme/digibyte-rosetta-node/configuration"
)

const (
	Transfer      = "transfer"
	Execution     = "execution"
	StatusSuccess = "success"
	StatusFail    = "fail"
	ActionTypeFee = "fee"
)

// DigibyteClientRPC is an implementation of DigibyteClient using RPC.
type DigibyteClientRPC struct {
	rpcConnConfig     *rpcclient.ConnConfig
	applicationConfig *configuration.Config
}

// NewDigibyteClient returns an implementation of DigibyteClient
func NewDigibyteClient(applicationConfig *configuration.Config) (cli DigibyteClient) {
	rpcConnConfig := rpcclient.ConnConfig{
		Host:         applicationConfig.Node.Endpoint,
		User:         applicationConfig.Node.Username,
		Pass:         applicationConfig.Node.Password,
		HTTPPostMode: true,                               // Bitcoin core only supports HTTP POST mode
		DisableTLS:   !applicationConfig.Node.TLSEnabled, // Bitcoin core does not provide TLS by default
	}

	return &DigibyteClientRPC{
		rpcConnConfig:     &rpcConnConfig,
		applicationConfig: applicationConfig,
	}
}

func (rpcClient *DigibyteClientRPC) reconnect() (client *rpcclient.Client) {
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(rpcClient.rpcConnConfig, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

// GetConfig retrieves the general application config that has been configured
func (rpcClient *DigibyteClientRPC) GetConfig() *configuration.Config {
	return rpcClient.applicationConfig
}

// GetStatus will return the Blockchain base info based on that node
func (rpcClient *DigibyteClientRPC) GetStatus(ctx context.Context) (*btcjson.GetBlockChainInfoResult, error) {
	client := rpcClient.reconnect()
	defer client.Shutdown()

	result, err := client.GetBlockChainInfo()
	return result, err
}

// GetBlock will return you the block specification for a given height
func (rpcClient *DigibyteClientRPC) GetBlock(ctx context.Context, height int64) (*wire.MsgBlock, error) {
	client := rpcClient.reconnect()
	defer client.Shutdown()

	result, err := client.GetBlockHash(height)
	block, err := client.GetBlock(result)
	return block, err
}

// GetBlockByHash will return you the block specification for a given height
func (rpcClient *DigibyteClientRPC) GetBlockByHash(ctx context.Context, hash string) (*wire.MsgBlock, error) {
	client := rpcClient.reconnect()
	defer client.Shutdown()

	blockHash, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		return nil, err
	}

	block, err := client.GetBlock(blockHash)
	if err != nil {
		return nil, err
	}

	return block, nil
}

// GetLatestBlock returns the latest Digibyte block.
func (rpcClient *DigibyteClientRPC) GetLatestBlock(ctx context.Context) (*wire.MsgBlock, error) {
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
