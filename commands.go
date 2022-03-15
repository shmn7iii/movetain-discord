package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// var nftAddress = "BKxCq9Q9nezwuQZvoe6oubcyr4w8F9NWgsAugpY71And"

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "mint",
			Description: "Mint a NFT",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "image",
					Description: "Image URL",
					Required:    true,
				},
			},
		},
		{
			Name:        "memo",
			Description: "Create memo transaction",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "content",
					Description: "Memmo content",
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"mint": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			imageURL := i.ApplicationCommandData().Options[0].StringValue()
			jsonCID := uploadJson2ipfs(imageURL)
			nftAddress, sig := mint(jsonCID)

			// if arg[0] == ipfs://
			// else if http:// or https://
			// else error

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: ":tada: Success! BOT mints a NFT." +
						"\n```" +
						"\nNFT Address:       " + nftAddress +
						"\nNFT Signature:     " + sig +
						"\nJSON IPFS Address: " + jsonCID +
						"```",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Emoji: discordgo.ComponentEmoji{
										Name: "🔍",
									},
									Label: "Check your NFT with Solana Explorer",
									Style: discordgo.LinkButton,
									URL: fmt.Sprintf(
										"https://explorer.solana.com/address/%s?cluster=devnet", nftAddress,
									),
								},
							},
						},
					},
				},
			})
			if err != nil {
				panic(err)
			}
		},
		"memo": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			content := i.ApplicationCommandData().Options[0].StringValue()
			txAddress := memo(content)

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: ":tada: Success! BOT Create a memo transaction." +
						"\n```" +
						"\nTransaction Address: " + txAddress +
						"```",
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Emoji: discordgo.ComponentEmoji{
										Name: "🔍",
									},
									Label: "Check your Transaction with Solana Explorer",
									Style: discordgo.LinkButton,
									URL: fmt.Sprintf(
										"https://explorer.solana.com/tx/%s?cluster=devnet", txAddress,
									),
								},
							},
						},
					},
				},
			})
			if err != nil {
				panic(err)
			}
		},
	}
)
