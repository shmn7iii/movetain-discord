# solana-go-sample
Solana no devnet de token wo mint suru discord bot.

## usage
```bash
$ go build
$ ./solana-go-sample -feepayer xxx -token yyy -guild zzz
```

- feepayer  
    FeePayer's base58
- token  
    Discord BOT's token
- guild  
    Guild ID for guild commands


## about

0. You need host IPFS Node.
1. On Discord, type this.
    ```
   /mint imageURL:<image url>
   ```
2. Then, BOT uploads JSON which include image URL to IPFS.
3. BOT mints a NFT.
4. Complete. You can check your NFT on Solana Explorer, and so on.