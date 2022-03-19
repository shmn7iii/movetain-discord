package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/rpc"
	"github.com/portto/solana-go-sdk/types"
)

var (
	JSON_KEYS      jsonKeys
	SOLANA_CLIENT  *client.Client
	DISCORD_CLIENT *discordgo.Session
	FEE_PAYER      types.Account

	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

type jsonKeys struct {
	BotToken       string `json:"botToken"`
	GuildID        string `json:"guildId"`
	FeePayerBase58 string `json:"feePayerBase58"`
}

func init() {
	flag.Parse()

	// Json読み込み
	raw, err := ioutil.ReadFile("secrets/keys.json")
	if err != nil {
		log.Fatalf("[Twitter] can't read secrets/keys.json: %v", err)
		return
	}
	json.Unmarshal(raw, &JSON_KEYS)

	FEE_PAYER, _ = types.AccountFromBase58(JSON_KEYS.FeePayerBase58)
}

// discord session
func init() {
	// open discord session
	var err error
	DISCORD_CLIENT, err = discordgo.New("Bot " + JSON_KEYS.BotToken)
	if err != nil {
		log.Fatalf("[Discord] Invalid bot parameters: %v", err)
	}
	// add command handlers
	DISCORD_CLIENT.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	// solana devnet
	SOLANA_CLIENT = client.NewClient(rpc.DevnetRPCEndpoint)
	resp, err := SOLANA_CLIENT.GetVersion(context.TODO())
	if err != nil {
		log.Fatalf("[Solana ] Failed to version info, err: %v", err)
	}
	log.Println("[Solana ] 🎉 Solana client has launched. version", resp.SolanaCore)

	// discord
	DISCORD_CLIENT.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("[Discord] 🥳 Logged in as", s.State.User.Username)
	})
	err = DISCORD_CLIENT.Open()
	if err != nil {
		log.Fatalf("[Discord] Cannot open the session: %v", err)
	}
	// add commands
	log.Println("[Discord] 🔨 Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := DISCORD_CLIENT.ApplicationCommandCreate(DISCORD_CLIENT.State.User.ID, JSON_KEYS.GuildID, v)
		if err != nil {
			log.Panicf("[Discord] Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	// defer close
	defer DISCORD_CLIENT.Close()
	// stop
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	fmt.Println("[Discord] 👋 bye bye.")
	// remove commands
	if *RemoveCommands {
		log.Println("[Discord] Removing commands...")
		for _, v := range registeredCommands {
			err := DISCORD_CLIENT.ApplicationCommandDelete(DISCORD_CLIENT.State.User.ID, JSON_KEYS.GuildID, v.ID)
			if err != nil {
				log.Panicf("[Discord] Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}
}
