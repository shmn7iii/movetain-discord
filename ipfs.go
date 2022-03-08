package main

import (
	"log"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

var ipfsShell *shell.Shell

func init() {
	ipfsShell = shell.NewShell("localhost:5001")
}

// Feature: 画像をIPFSにあげるかどうかは未定
//func uploadImage2ipfs(imageURL string) (imageCID string) {
//	imageCID, err := ipfsShell.Add(
//		strings.NewReader(
//			"{" +
//				"\n  \"name\": \"Super Test NFT\"," +
//				"\n  \"description\": \"Metyakutya Tekitou na NFT no Test.\"," +
//				"\n  \"image\": \"" + imageURL + "\"," +
//				"\n  \"external_url\": \"https://www.github.com/shmn7iii/solana-go-example\"" +
//				"\n}" +
//				""),
//	)
//	if err != nil {
//		log.Fatalf("failed to add file to ipfs, err: %v", err)
//	}
//	return
//}

func uploadJson2ipfs(imageURI string) (jsonCID string) {
	jsonCID, err := ipfsShell.Add(
		strings.NewReader(
			"{" +
				"\n  \"name\": \"Super Test NFT\"," +
				"\n  \"description\": \"Metyakutya Tekitou na NFT no Test.\"," +
				"\n  \"image\": \"" + imageURI + "\"," +
				"\n  \"external_url\": \"https://www.github.com/shmn7iii/solana-go-example\"" +
				"\n}" +
				""),
	)
	if err != nil {
		log.Fatalf("failed to add file to ipfs, err: %v", err)
	}
	return
}
