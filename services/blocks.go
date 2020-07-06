package services

import (
	"context"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/marpme/digibyte-rosetta-node/client"
)

// BlockAPIService client based implementation of the block servicer
type blockAPIService struct {
	server.BlockAPIServicer
	client client.DigibyteClient
}

// NewBlockAPIService creates a new block API service
func NewBlockAPIService(client client.DigibyteClient) server.BlockAPIServicer {
	return &blockAPIService{
		client: client,
	}
}

func mapTransactions(txs []string) []*types.TransactionIdentifier {
	var transactionsIdentifiers []*types.TransactionIdentifier
	for i := 0; i < len(txs); i++ {
		transactionsIdentifiers = append(transactionsIdentifiers, &types.TransactionIdentifier{
			Hash: txs[i],
		})
	}

	return transactionsIdentifiers
}

func (blockService *blockAPIService) retriveBlock(ctx context.Context, blockRequest *types.BlockRequest) (*btcjson.GetBlockVerboseResult, *btcjson.GetBlockVerboseResult, *types.Error) {

	var block, prevBlock *btcjson.GetBlockVerboseResult
	var err error

	if blockRequest.BlockIdentifier.Index != nil {
		block, err = blockService.client.GetBlock(ctx, *blockRequest.BlockIdentifier.Index)
		prevBlock, err = blockService.client.GetBlock(ctx, *blockRequest.BlockIdentifier.Index-1)
	} else if blockRequest.BlockIdentifier.Hash != nil {
		block, err = blockService.client.GetBlockByHash(ctx, *blockRequest.BlockIdentifier.Hash)

		if err != nil {
			return nil, nil, ErrUnableToGetBlk
		}
		prevBlock, err = blockService.client.GetBlockByHash(ctx, block.PreviousHash)
	} else {
		block, err := blockService.client.GetLatestBlock(ctx)
		if err != nil {
			return nil, nil, ErrUnableToGetBlk
		}
		prevBlock, err = blockService.client.GetBlockByHash(ctx, block.Header.PrevBlock.String())
	}

	if err != nil {
		return nil, nil, ErrUnableToGetBlk
	}

	return block, prevBlock, nil
}

// Block retrieves the block for a given candidate
func (blockService *blockAPIService) Block(ctx context.Context, blockRequest *types.BlockRequest) (*types.BlockResponse, *types.Error) {
	block, prevBlock, err := blockService.retriveBlock(ctx, blockRequest)

	if err != nil {
		return nil, err
	}

	return &types.BlockResponse{
		Block: &types.Block{
			BlockIdentifier: &types.BlockIdentifier{
				Hash: block.Hash,
			},
			ParentBlockIdentifier: &types.BlockIdentifier{
				Hash: prevBlock.Hash,
			},
			Timestamp: block.Time,
		},
		OtherTransactions: mapTransactions(block.Tx),
	}, nil
}

// BlockTransaction retrieves the block with the given transactions included
func (blockService *blockAPIService) BlockTransaction(ctx context.Context, blockTransaction *types.BlockTransactionRequest) (*types.BlockTransactionResponse, *types.Error) {
	block, err := blockService.client.GetBlockByHashWithTransaction(ctx, blockTransaction.BlockIdentifier.Hash)

	if err != nil {
		return nil, ErrUnableToGetBlk
	}

	for index, tx := range block.Transactions {
		if tx.TxHash().String() == blockTransaction.TransactionIdentifier.Hash {
			txOperations := make([]*types.Operation, len(tx.TxOut))

			for _, vOut := range block.Transactions[index].TxOut {
				txOperations.append(txOperations, &types.Operation{
					OperationIdentifier: &types.OperationIdentifier {
						Index: vOut.PkScript
					}
				})
			}

			return &types.BlockTransactionResponse{
				Transaction: &types.Transaction{
					TransactionIdentifier: &types.TransactionIdentifier{
						Hash: tx.TxHash().String(),
					},
					Metadata: map[string]interface{}{
						"size":     tx.SerializeSize(),
						"lockTime": tx.LockTime,
					},
					Operations: txOperations,
				},
			}, nil
		}
	}

	return nil, ErrUnableToGetTxns
}
