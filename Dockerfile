FROM golang:1.15

ENV CGO_ENABLED=0 \
    GOOS=linux \
	GOARCH=amd64 \
    GO111MODULE=on

RUN mkdir -p /go/src/github.com/helixauth/helix

WORKDIR /go/src/github.com/helixauth/helix
COPY . .

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -log-prefix=false --build="go build ./src/main.go" --command=./main
