package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/program/sysprog"
	"github.com/portto/solana-go-sdk/program/tokenprog"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

// cli flag
var (
	FeePayerKeyPairBase58 = flag.String("feepayer", "", "keypair no base 58  *private key no base 58 dato error deru")
	AliceKeyPairBase58    = flag.String("alice", "", "alice no base 58")
)

// [global] solana client
var solanaClient *client.Client

// [global] account
var feePayerAccount types.Account
var aliceAccount types.Account

func loadFeePayer() {
	flag.Parse()
	fe, err := types.AccountFromBase58(*FeePayerKeyPairBase58)
	if err != nil {
		log.Fatalf("load fee payer, err: %v", err)
	}
	bal, err := solanaClient.GetBalance(context.TODO(), fe.PublicKey.ToBase58())
	if err != nil {
		log.Fatalf("load fee payer's balance, err: %v", err)
	}
	fmt.Println("💰 Fee payer:", fe.PublicKey.ToBase58())
	fmt.Println("     balance:", bal)
	feePayerAccount = fe

	al, err := types.AccountFromBase58(*AliceKeyPairBase58)
	if err != nil {
		log.Fatalf("load fee payer, err: %v", err)
	}
	bal, err = solanaClient.GetBalance(context.TODO(), al.PublicKey.ToBase58())
	if err != nil {
		log.Fatalf("load fee payer's balance, err: %v", err)
	}
	fmt.Println("🙋‍♀️ Alice:", al.PublicKey.ToBase58())
	fmt.Println("  balance:", bal)
	aliceAccount = al
}

type TokenMintResponse struct {
	Tx string `json:"Tx"`
}

func mint(c echo.Context) error {

	// get init balance
	rentExemptionBalance, err := solanaClient.GetMinimumBalanceForRentExemption(
		context.Background(),
		tokenprog.MintAccountSize,
	)
	if err != nil {
		log.Fatalf("get min balacne for rent exemption, err: %v", err)
	}

	// create accounts
	mintAccount := types.NewAccount()
	fmt.Println("mint:", mintAccount.PublicKey.ToBase58())

	// //air drop
	// airdrop_txhash, err := solanaClient.RequestAirdrop(
	// 	context.Background(),
	// 	feePayerAccount.PublicKey.ToBase58(),
	// 	1e9, // 1 SOL = 10^9 lamports
	// )
	// if err != nil {
	// 	log.Fatalf("air drop error, err: %v\n", err)
	// }

	// get blockhash
	res, err := solanaClient.GetRecentBlockhash(context.Background())
	if err != nil {
		log.Fatalf("get recent block hash error, err: %v\n", err)
	}

	// create transaction
	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: types.NewMessage(types.NewMessageParam{
			FeePayer:        mintAccount.PublicKey,
			RecentBlockhash: res.Blockhash,
			Instructions: []types.Instruction{
				sysprog.CreateAccount(sysprog.CreateAccountParam{
					From:     feePayerAccount.PublicKey,
					New:      mintAccount.PublicKey,
					Owner:    common.TokenProgramID,
					Lamports: rentExemptionBalance,
					Space:    tokenprog.MintAccountSize,
				}),
				tokenprog.InitializeMint(tokenprog.InitializeMintParam{
					Decimals:   8,
					Mint:       mintAccount.PublicKey,
					MintAuth:   aliceAccount.PublicKey,
					FreezeAuth: nil,
				}),
			},
		}),
		Signers: []types.Account{feePayerAccount, mintAccount},
	})
	if err != nil {
		log.Fatalf("generate tx error, err: %v\n", err)
	}

	txhash, err := solanaClient.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("send tx error, err: %v\n", err)
	}

	response := new(TokenMintResponse)
	response.Tx = txhash
	return c.JSON(http.StatusOK, res)
}

func main() {
	// solana devnet
	solanaClient = client.NewClient(rpc.DevnetRPCEndpoint)
	resp, err := solanaClient.GetVersion(context.TODO())
	if err != nil {
		log.Fatalf("failed to version info, err: %v", err)
	}
	fmt.Println("\n🎉 Solana client has launched. version", resp.SolanaCore)
	// load fee payer
	loadFeePayer()

	// echo
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/mint", mint)
	e.Logger.Fatal(e.Start(":1323"))

}
