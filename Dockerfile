# # Stage 0, pull everything
# FROM alpine/git as cloner
# WORKDIR /app
# RUN git clone https://github.com/antigravities/levelup.git

# Stage 1
FROM node:10-alpine as builder-webpack
RUN npm install webpack webpack-cli -g
WORKDIR /app
COPY ./assets/ /app/
RUN npm install
RUN webpack

# GO111MODULE=on go build -i
# Stage 2, build Go
FROM golang:1.14.3-alpine AS builder-go
WORKDIR /src
COPY . .
RUN GO111MODULE=on go build -i
 

# Stage 3, run the app
FROM alpine:latest AS server
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=builder-webpack /static/ ./static/
COPY --from=builder-go /src/levelup .

EXPOSE 4000

CMD ["./levelup"]
