package main

import (
	"context"
	"log"

	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/pkg/pointer"
	"github.com/portto/solana-go-sdk/program/assotokenprog"
	"github.com/portto/solana-go-sdk/program/metaplex/tokenmeta"
	"github.com/portto/solana-go-sdk/program/sysprog"
	"github.com/portto/solana-go-sdk/program/tokenprog"
	"github.com/portto/solana-go-sdk/types"
)

func mint(jsonCID string) (nftAddress, sig string) {
	mint := types.NewAccount()
	nftAddress = mint.PublicKey.ToBase58()

	ata, _, err := common.FindAssociatedTokenAddress(FEE_PAYER.PublicKey, mint.PublicKey)
	if err != nil {
		log.Fatalf("failed to find a valid ata, err: %v", err)
	}

	tokenMetadataPubkey, err := tokenmeta.GetTokenMetaPubkey(mint.PublicKey)
	if err != nil {
		log.Fatalf("failed to find a valid token metadata, err: %v", err)

	}

	tokenMasterEditionPubkey, err := tokenmeta.GetMasterEdition(mint.PublicKey)
	if err != nil {
		log.Fatalf("failed to find a valid master edition, err: %v", err)
	}

	mintAccountRent, err := SOLANA_CLIENT.GetMinimumBalanceForRentExemption(context.Background(), tokenprog.MintAccountSize)
	if err != nil {
		log.Fatalf("failed to get mint account rent, err: %v", err)
	}

	recentBlockhashResponse, err := SOLANA_CLIENT.GetRecentBlockhash(context.Background())
	if err != nil {
		log.Fatalf("failed to get recent blockhash, err: %v", err)
	}

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Signers: []types.Account{mint, FEE_PAYER},
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        FEE_PAYER.PublicKey,
			RecentBlockhash: recentBlockhashResponse.Blockhash,
			Instructions: []types.Instruction{
				sysprog.CreateAccount(sysprog.CreateAccountParam{
					From:     FEE_PAYER.PublicKey,
					New:      mint.PublicKey,
					Owner:    common.TokenProgramID,
					Lamports: mintAccountRent,
					Space:    tokenprog.MintAccountSize,
				}),
				tokenprog.InitializeMint(tokenprog.InitializeMintParam{
					Decimals: 0,
					Mint:     mint.PublicKey,
					MintAuth: FEE_PAYER.PublicKey,
				}),
				tokenmeta.CreateMetadataAccount(tokenmeta.CreateMetadataAccountParam{
					Metadata:                tokenMetadataPubkey,
					Mint:                    mint.PublicKey,
					MintAuthority:           FEE_PAYER.PublicKey,
					Payer:                   FEE_PAYER.PublicKey,
					UpdateAuthority:         FEE_PAYER.PublicKey,
					UpdateAuthorityIsSigner: true,
					IsMutable:               true,
					MintData: tokenmeta.Data{
						Name:                 "Super Test NFT",
						Symbol:               "STT",
						Uri:                  "https://gateway.ipfs.io/ipfs/" + jsonCID,
						SellerFeeBasisPoints: 100,
						Creators: &[]tokenmeta.Creator{
							{
								Address:  FEE_PAYER.PublicKey,
								Verified: true,
								Share:    100,
							},
						},
					},
				}),
				assotokenprog.CreateAssociatedTokenAccount(assotokenprog.CreateAssociatedTokenAccountParam{
					Funder:                 FEE_PAYER.PublicKey,
					Owner:                  FEE_PAYER.PublicKey,
					Mint:                   mint.PublicKey,
					AssociatedTokenAccount: ata,
				}),
				tokenprog.MintTo(tokenprog.MintToParam{
					Mint:   mint.PublicKey,
					To:     ata,
					Auth:   FEE_PAYER.PublicKey,
					Amount: 1,
				}),
				tokenmeta.CreateMasterEdition(tokenmeta.CreateMasterEditionParam{
					Edition:         tokenMasterEditionPubkey,
					Mint:            mint.PublicKey,
					UpdateAuthority: FEE_PAYER.PublicKey,
					MintAuthority:   FEE_PAYER.PublicKey,
					Metadata:        tokenMetadataPubkey,
					Payer:           FEE_PAYER.PublicKey,
					MaxSupply:       pointer.Uint64(0),
				}),
			},
		}),
	})
	if err != nil {
		log.Fatalf("failed to new a tx, err: %v", err)
	}

	sig, err = SOLANA_CLIENT.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("failed to send tx, err: %v", err)
	}

	log.Println("[Solana ] 🪪 BOT has minted a NFT")
	log.Println("[Solana ]      Account:  ", nftAddress)
	log.Println("[Solana ]      Signature:", sig)

	return
}
