#!/bin/bash

# webpack
cd assets
npm i
npm run webpack --mode="production"

# build app
go build -i
graftcp ./levelup