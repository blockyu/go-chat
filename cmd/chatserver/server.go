package chatserver

import "flag"

import (
	"github.com/blockyu/go-chat"
)

func main() {
	tcpAddr := flag.String("tcp", "localhost:8989", "Address for the TCP chat server to listen on")
	wsAddr := flag.String("ws", "localhost:8099", "Address for websocket chat server to listen on")
	flag.Parse()
	api := chatapi.New()
}
