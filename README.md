# movetain-discord
Movetain kun no Dicord BOT version. Solana no test no tameni tukutta yo. 

## usage

First, set `secrets/keys.json`

```json
{
  "botToken": "<Discord BOT Token>",
  "guildId": "<Discord Guild ID>",
  "feePayerBase58": "<Base58 of FeePayer's keypair>"
}
```

Then,

```bash
$ go build -o bin/main
$ ./bin/main
```

## about

0. You need host IPFS Node.
1. On Discord, type this.
   ```
   /mint imageURL:<image url>
   ```
   
   ```
   /memo content:<memo content>
   ```
2. Then, BOT uploads JSON which include image URL to IPFS.
3. BOT mints a NFT.
4. Complete. You can check your NFT on Solana Explorer, and so on.

＊ **YOUR METAPLEX JSON is stored on IPFS, but YOUR IMAGE is not.**
