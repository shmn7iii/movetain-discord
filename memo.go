package main

import (
	"context"
	"log"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/memoprog"
	"github.com/portto/solana-go-sdk/types"
)

func memo(content string) string {
	// fetch recent blockhash
	recentBlockhashResponse, err := solanaClient.GetRecentBlockhash(context.Background())
	if err != nil {
		log.Fatalf("failed to get recent blockhash, err: %v", err)
	}

	// create a tx
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{feePayer, feePayer},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        feePayer.PublicKey,
			RecentBlockhash: recentBlockhashResponse.Blockhash,
			Instructions: []types.Instruction{
				// memo instruction
				memoprog.BuildMemo(memoprog.BuildMemoParam{
					SignerPubkeys: []common.PublicKey{feePayer.PublicKey},
					Memo:          []byte(content),
				}),
			},
		}),
	})
	if err != nil {
		log.Fatalf("failed to new a transaction, err: %v", err)
	}

	// send tx
	txhash, err := solanaClient.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("failed to send tx, err: %v", err)
	}

	log.Println("[Solana ] 🧱 BOT has created a transaction")
	log.Println("[Solana ]      Tx Hash:", txhash)

	return txhash
}
