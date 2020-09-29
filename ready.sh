#!/bin/bash

# webpack
cd assets
npm i
npm run webpack --mode="production"

# build app
cd ..
go build -i
graftcp ./levelup