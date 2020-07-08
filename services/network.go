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

// NewNetworkAPIService creates a new service to communicate about Network related topics
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
			RosettaVersion: cfg.Version.RosettaVersion,
			NodeVersion:    cfg.Version.DigibyteVersion,
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
				client.Transfer,
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
	terr := ValidateNetworkIdentifier(ctx, network.client, request.NetworkIdentifier)
	if terr != nil {
		return nil, terr
	}

	status, err := network.client.GetStatus(ctx)
	if err != nil {
		return nil, ErrUnableToGetNodeStatus
	}

	height := int64(status.Blocks)

	bestBlock, err := network.client.GetBlock(ctx, height)
	if err != nil {
		return nil, ErrUnableToGetNodeStatus
	}

	genesisBlock, err := network.client.GetBlock(ctx, 0)
	if err != nil {
		return nil, ErrUnableToGetNodeStatus
	}

	resp := &types.NetworkStatusResponse{
		CurrentBlockIdentifier: &types.BlockIdentifier{
			Index: height,
			Hash:  bestBlock.Hash,
		},
		CurrentBlockTimestamp: bestBlock.Time * 1000, // ms
		GenesisBlockIdentifier: &types.BlockIdentifier{
			Index: 0,
			Hash:  genesisBlock.Hash,
		},
		Peers: nil,
	}

	return resp, nil
}
