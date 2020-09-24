# Level Up

Level Up is a Steam game recommendations site where users can submit their own ideas. It was originally designed for the /r/Steam Discord server.

## Try it

/r/Steam's copy of Level Up is available at https://recommendations.steamsal.es/.

## Run it yourself

In order to work as intended, Level Up requires a US IP address. If you don't have one on your machine, use [graftcp](https://github.com/hmgle/graftcp) along with a socks5 proxy (normally proxychains, etc would work but these rely on dynamically linked libraries which Go doesn't do)

## Hack it

Level Up's server components are written in Go and operate on top of an Amazon DynamoDB database.