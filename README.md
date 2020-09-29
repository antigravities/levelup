# Level Up

Level Up is a Steam game recommendations site where users can submit their own ideas. It was originally designed for the /r/Steam Discord server.

## Try it

/r/Steam's copy of Level Up is available at https://recommendations.steamsal.es/.

## Run it yourself

Level Up's server components are written in Go and operate on top of an Amazon DynamoDB database. You will need both - AWS's free tier contains 1 free DynamoDB table.

In order to work as intended, Level Up requires a US IP address. If you don't have one on your machine, use [graftcp](https://github.com/hmgle/graftcp) along with a socks5 proxy (normally proxychains, etc would work but these rely on dynamically linked libraries which Go doesn't do).

1. Download and install Go 1.14+.
2. Download and install Node.js + npm.
3. `cp .env.example .env && $EDITOR .env`
3. Start: 
* If you are not using `graftcp`: `cd assets && npm run webpack && cd .. && go build -i && ./levelup`
* If you are using `graftcp`: `./ready.sh` will do all of the above for you