package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

var (
	FeePayerBase58 = flag.String("feepayer", "", "FeePayer no base58. "+
		"*keypair no base 58, private key no base 58 dato error deru")
	BotToken       = flag.String("token", "", "Bot access token")
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var solanaClient *client.Client
var discordSession *discordgo.Session

var feePayer types.Account

func init() {
	flag.Parse()
	feePayer, _ = types.AccountFromBase58(*FeePayerBase58)
}

// discord session
func init() {
	// open discord session
	var err error
	discordSession, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("[Discord] Invalid bot parameters: %v", err)
	}
	// add command handlers
	discordSession.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	// solana devnet
	solanaClient = client.NewClient(rpc.DevnetRPCEndpoint)
	resp, err := solanaClient.GetVersion(context.TODO())
	if err != nil {
		log.Fatalf("[Solana ] Failed to version info, err: %v", err)
	}
	log.Println("[Solana ] 🎉 Solana client has launched. version", resp.SolanaCore)

	// discord
	discordSession.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("[Discord] 🥳 Logged in as", s.State.User.Username)
	})
	err = discordSession.Open()
	if err != nil {
		log.Fatalf("[Discord] Cannot open the session: %v", err)
	}
	// add commands
	log.Println("[Discord] 🔨 Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := discordSession.ApplicationCommandCreate(discordSession.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("[Discord] Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	// defer close
	defer discordSession.Close()
	// stop
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	fmt.Println("[Discord] 👋 bye bye.")
	// remove commands
	if *RemoveCommands {
		log.Println("[Discord] Removing commands...")
		for _, v := range registeredCommands {
			err := discordSession.ApplicationCommandDelete(discordSession.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("[Discord] Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}
}
