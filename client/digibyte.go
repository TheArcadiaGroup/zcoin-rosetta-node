// Copyright (c) 2020 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package client

import (
	"crypto/tls"
	"sync"

	"github.com/marpme/digibyte-rosetta-node/configuration"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
)

const (
	Transfer      = "transfer"
	Execution     = "execution"
	StatusSuccess = "success"
	StatusFail    = "fail"
	ActionTypeFee = "fee"
)

type (
	// DigibyteClient is the Digibyte blockchain client interface.
	DigibyteClient interface {
		// 	// GetChainID returns the network chain context, derived from the
		// 	// genesis document.
		// 	GetChainID(ctx context.Context) (string, error)

		// 	// GetBlock returns the IoTex block at given height.
		// 	GetBlock(ctx context.Context, height int64) (*IoTexBlock, error)

		// 	// GetLatestBlock returns latest IoTex block.
		// 	GetLatestBlock(ctx context.Context) (*IoTexBlock, error)

		// 	// GetGenesisBlock returns the IoTex genesis block.
		// 	GetGenesisBlock(ctx context.Context) (*IoTexBlock, error)

		// 	// GetAccount returns the IoTex staking account for given owner address
		// 	// at given height.
		// 	GetAccount(ctx context.Context, height int64, owner string) (*Account, error)

		// 	// SubmitTx submits the given encoded transaction to the node.
		// 	SubmitTx(ctx context.Context, tx *iotextypes.Action) (txid string, err error)

		// // GetStatus returns the status overview of the node.
		// GetStatus(ctx context.Context) (*iotexapi.GetChainMetaResponse, error)

		// 	// GetVersion returns the server's version.
		// 	GetVersion(ctx context.Context) (*iotexapi.GetServerMetaResponse, error)

		// 	// GetTransactions returns transactions of the block.
		// 	GetTransactions(ctx context.Context, height int64) ([]*types.Transaction, error)

		// GetConfig returns the config.
		GetConfig() *configuration.Config
	}

	// IoTexBlock is the IoTex blockchain's block.
	DigibyteBlock struct {
		Height       int64  // Block height.
		Hash         string // Block hash.
		Timestamp    int64  // UNIX time, converted to milliseconds.
		ParentHeight int64  // Height of parent block.
		ParentHash   string // Hash of parent block.
	}

	Account struct {
		Nonce   uint64
		Balance string
	}

	// grpcDigibyteClient is an implementation of DigibyteClient using gRPC.
	grpcDigibyteClient struct {
		sync.RWMutex

		endpoint string
		grpcConn *grpc.ClientConn
		cfg      *configuration.Config
	}
)

// NewDigibyteClient returns an implementation of DigibyteClient
func NewDigibyteClient(cfg *configuration.Config) (cli DigibyteClient, err error) {
	grpc, err := grpc.Dial(cfg.Server.Endpoint, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	if err != nil {
		return
	}
	cli = &grpcDigibyteClient{grpcConn: grpc, cfg: cfg}
	return
}

func (c *grpcDigibyteClient) reconnect() (err error) {
	c.Lock()
	defer c.Unlock()
	// Check if the existing connection is good.
	if c.grpcConn != nil && c.grpcConn.GetState() != connectivity.Shutdown {
		return
	}
	c.grpcConn, err = grpc.Dial(c.endpoint, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	return err
}

func (c *grpcDigibyteClient) GetConfig() *configuration.Config {
	return c.cfg
}
