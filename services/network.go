// Copyright (c) 2020 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package services

import (
	"context"

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/marpme/digibyte-rosetta-node/client"
)

type networkAPIService struct {
	client client.DigibyteClient
}

func NewNetworkAPIService(client client.DigibyteClient) server.NetworkAPIServicer {
	return &networkAPIService{
		client: client,
	}
}

func (network *networkAPIService) NetworkList(context.Context, *types.MetadataRequest) (*types.NetworkListResponse, *types.Error) {
	cfg := network.client.GetConfig()
	return &types.NetworkListResponse{
		NetworkIdentifiers: []*types.NetworkIdentifier{{
			Blockchain: cfg.NetworkIdentifier.Blockchain,
			Network:    cfg.NetworkIdentifier.Network,
		},
		},
	}, nil
}

func (network *networkAPIService) NetworkOptions(
	ctx context.Context,
	request *types.NetworkRequest,
) (*types.NetworkOptionsResponse, *types.Error) {
	terr := ValidateNetworkIdentifier(ctx, network.client, request.NetworkIdentifier)
	if terr != nil {
		return nil, terr
	}
	cfg := network.client.GetConfig()
	return &types.NetworkOptionsResponse{
		Version: &types.Version{
			RosettaVersion: cfg.Server.RosettaVersion,
			NodeVersion:    cfg.Server.DigibyteVersion,
		},
		Allow: &types.Allow{
			OperationStatuses: []*types.OperationStatus{
				{
					Status:     client.StatusSuccess,
					Successful: true,
				},
				{
					Status:     client.StatusFail,
					Successful: false,
				},
			},
			OperationTypes: []string{
				client.ActionTypeFee,
				client.Transfer,
				client.Execution,
			},
			Errors: ErrorList,
		},
	}, nil
}

// ValidateNetworkIdentifier validates the network identifier.
func ValidateNetworkIdentifier(ctx context.Context, client client.DigibyteClient, ni *types.NetworkIdentifier) *types.Error {
	if ni != nil {
		cfg := client.GetConfig()
		if ni.Blockchain != cfg.NetworkIdentifier.Blockchain {
			return ErrInvalidBlockchain
		}
		if ni.SubNetworkIdentifier != nil {
			return ErrInvalidSubnetwork
		}
		if ni.Network != cfg.NetworkIdentifier.Network {
			return ErrInvalidNetwork
		}
	} else {
		return ErrMissingNID
	}
	return nil
}

func (network *networkAPIService) NetworkStatus(
	ctx context.Context,
	request *types.NetworkRequest,
) (*types.NetworkStatusResponse, *types.Error) {
	// terr := ValidateNetworkIdentifier(ctx, network.client, request.NetworkIdentifier)
	// if terr != nil {
	// 	return nil, terr
	// }

	// status, err := network.client.GetStatus(ctx)
	// if err != nil {
	// 	return nil, ErrUnableToGetNodeStatus
	// }
	// hei := int64(status.GetChainMeta().GetHeight())
	// blk, err := network.client.GetBlock(ctx, hei)
	// if err != nil {
	// 	return nil, ErrUnableToGetNodeStatus
	// }
	// genesisblk, err := network.client.GetBlock(ctx, 1)
	// if err != nil {
	// 	return nil, ErrUnableToGetNodeStatus
	// }
	// resp := &types.NetworkStatusResponse{
	// 	CurrentBlockIdentifier: &types.BlockIdentifier{
	// 		Index: hei,
	// 		Hash:  blk.Hash,
	// 	},
	// 	CurrentBlockTimestamp: blk.Timestamp, // ms
	// 	GenesisBlockIdentifier: &types.BlockIdentifier{
	// 		Index: genesisblk.Height,
	// 		Hash:  genesisblk.Hash,
	// 	},
	// 	Peers: nil,
	// }

	// return resp, nil
	return nil, nil
}
