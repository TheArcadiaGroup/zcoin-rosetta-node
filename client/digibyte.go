// Copyright (c) 2020 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package client

import (
	"context"
	"log"
	"sync"

	"github.com/btcsuite/btcd/btcjson"
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

type NetworkInfo struct {
	version string `json:"version"`
}

type (
	// DigibyteClient is the Digibyte blockchain client interface.
	DigibyteClient interface {

		// GetBlock returns the IoTex block at given height.
		GetBlock(ctx context.Context, height int64) (*wire.MsgBlock, error)

		// GetStatus returns the status overview of the node.
		GetStatus(ctx context.Context) (*btcjson.GetBlockChainInfoResult, error)

		// GetConfig returns the config.
		GetConfig() *configuration.Config
	}

	// DigibyteRPCClient is an implementation of DigibyteClient using RPC.
	DigibyteRPCClient struct {
		sync.RWMutex

		endpoint          string
		rpcConnConfig     *rpcclient.ConnConfig
		applicationConfig *configuration.Config
	}
)

// NewDigibyteClient returns an implementation of DigibyteClient
func NewDigibyteClient(cfg *configuration.Config) (cli DigibyteClient, err error) {
	rpcConnConfig := rpcclient.ConnConfig{
		Host:         cfg.Server.Endpoint,
		User:         cfg.Server.Username,
		Pass:         cfg.Server.Password,
		HTTPPostMode: true,                   // Bitcoin core only supports HTTP POST mode
		DisableTLS:   !cfg.Server.TLSEnabled, // Bitcoin core does not provide TLS by default
	}

	return &DigibyteRPCClient{
		rpcConnConfig:     &rpcConnConfig,
		applicationConfig: cfg,
	}, nil
}

func (rpcClient *DigibyteRPCClient) reconnect() (client *rpcclient.Client) {
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(rpcClient.rpcConnConfig, nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func (rpcClient *DigibyteRPCClient) GetConfig() *configuration.Config {
	return rpcClient.applicationConfig
}

func (rpcClient *DigibyteRPCClient) GetStatus(ctx context.Context) (*btcjson.GetBlockChainInfoResult, error) {
	client := rpcClient.reconnect()
	defer client.Shutdown()

	result, err := client.GetBlockChainInfo()
	return result, err
}

func (rpcClient *DigibyteRPCClient) GetBlock(ctx context.Context, height int64) (*wire.MsgBlock, error) {
	client := rpcClient.reconnect()
	defer client.Shutdown()

	result, err := client.GetBlockHash(height)
	block, err := client.GetBlock(result)
	return block, err
}
