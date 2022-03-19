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
	recentBlockhashResponse, err := SOLANA_CLIENT.GetRecentBlockhash(context.Background())
	if err != nil {
		log.Fatalf("failed to get recent blockhash, err: %v", err)
	}

	// create a tx
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{FEE_PAYER, FEE_PAYER},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        FEE_PAYER.PublicKey,
			RecentBlockhash: recentBlockhashResponse.Blockhash,
			Instructions: []types.Instruction{
				// memo instruction
				memoprog.BuildMemo(memoprog.BuildMemoParam{
					SignerPubkeys: []common.PublicKey{FEE_PAYER.PublicKey},
					Memo:          []byte(content),
				}),
			},
		}),
	})
	if err != nil {
		log.Fatalf("failed to new a transaction, err: %v", err)
	}

	// send tx
	txhash, err := SOLANA_CLIENT.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("failed to send tx, err: %v", err)
	}

	log.Println("[Solana ] 🧱 BOT has created a transaction")
	log.Println("[Solana ]      Tx Hash:", txhash)

	return txhash
}
