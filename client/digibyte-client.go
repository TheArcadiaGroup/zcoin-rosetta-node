package client

import (
	"context"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/wire"
	"github.com/marpme/digibyte-rosetta-node/configuration"
)

// DigibyteClient is the Digibyte blockchain client interface
type DigibyteClient interface {
	// GetBlock returns the IoTex block at given height.
	GetBlock(ctx context.Context, height int64) (*wire.MsgBlock, error)

	// GetStatus returns the status overview of the node.
	GetStatus(ctx context.Context) (*btcjson.GetBlockChainInfoResult, error)

	// GetConfig returns the config.
	GetConfig() *configuration.Config
}
