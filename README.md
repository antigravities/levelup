# Level Up

Level Up is a Steam game recommendations site where users can submit their own apps. It was originally designed for the /r/Steam Discord server.

## Try it

/r/Steam's copy of Level Up is available at https://recommendations.steamsal.es/.

## Run it yourself

### 1. using Docker

Make sure you have the .env file if not run: `cp .env.example .env && $EDITOR .env`

`docker build -t levelup .` 

`docker run --env-file .env -p 4000:4000 levelup` 

or

`docker run --env-file .env -p 4000:4000 bikerwhocodes/levelup` 

or

`docker-compose up`

or 

`docker-compose up -d` for headless mode

### 2. Natively

Level Up's server components are written in Go and operate on top of an Amazon DynamoDB database (DynamoDB's free tier is more than enough to host).

1. Download and install Go 1.14+.
2. Download and install Node.js + npm.
3. `cp .env.example .env && $EDITOR .env`
3. Start: `./ready.sh`, or
```sh
cd assets
webpack
cd ..

go build -i
./levelup
```

**Important:** In order to work as intended, Level Up unfortunately requires a US IP address. If you don't have one on your machine, launch Level Up in `fetch` mode using a US IP by running `./ready.sh fetch` (this will not start a Web server - launch levelup in `serve` mode similarly to the above to do that).

## Contributing

Contributions must be accompanied by a Signed-off-by header certifying your commit(s) under the [Developer Certificate of Origin](https://developercertificate.org/).

## License

```
Copyright (c) 2020 Alexandra Frock, Cutie Café, contributors

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
```
