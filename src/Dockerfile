FROM golang:1.19

WORKDIR /galibot

#TODO: Only copy go files
COPY . .

RUN go mod download

ENTRYPOINT go build && ./galibot
