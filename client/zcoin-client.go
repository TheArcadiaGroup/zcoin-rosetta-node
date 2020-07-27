package client

import (
	"context"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/wire"
	"gopkg.in/ArcadiaMediaGroup/zcoin-rosetta-node.v0/configuration"
)

// ZcoinClient is the Zcoin blockchain client interface
type ZcoinClient interface {
	// GetBlock returns the Zcoin block at given height.
	GetBlock(ctx context.Context, height int64) (*btcjson.GetBlockVerboseResult, error)

	// GetBlock returns the Zcoin block with a given hash.
	GetBlockByHash(ctx context.Context, hash string) (*btcjson.GetBlockVerboseResult, error)

	// GetBlock returns the Zcoin block with a given hash.
	GetBlockByHashWithTransaction(ctx context.Context, hash string) (*btcjson.GetBlockVerboseTxResult, error)

	// GetLatestBlock returns the latest Zcoin block.
	GetLatestBlock(ctx context.Context) (*wire.MsgBlock, error)

	// GetStatus returns the status overview of the node.
	GetStatus(ctx context.Context) (*btcjson.GetBlockChainInfoResult, error)

	// GetConfig returns the config.
	GetConfig() *configuration.Config
}
