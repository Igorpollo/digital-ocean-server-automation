FROM golang:alpine

ENV GO111MODULE=on

WORKDIR /app

COPY . .

COPY $HOME/.shh/ $HOME/.ssh/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

EXPOSE 5000
ENTRYPOINT ["/app/automation"]