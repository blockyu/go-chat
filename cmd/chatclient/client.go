package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	name := fmt.Sprintf("User%s", rand.Intn(400))
	clientName := flag.String("client", name, "name of the client to connect")
	url := flag.String("address", "ws://localhost:8099/mychat", "address of the chat server")
	roomName := flag.String("room", "lounge", "name off the chat room")
	origin := flag.String("origin", "http://localhost:8099", "origin flag for the conn client")
	flag.Parse()

	ws, err := websocket.Dial(*url, "", *origin)
	defer ws.Close()
	if err != nil {
		log.Fatal("conn dial error", err)
	}
	go func() {
		scanner := bufio.NewScanner(ws)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		log.Println("Connection seem to be cloased or error occured", scanner.Err())
		os.Exit(0)
	}()
	clientInfo := struct {
		Name string `json:name`
		Room string `json:room`
	}{*clientName, *roomName}
	err = websocket.JSON.Send(ws, &clientInfo)
	if err != nil {
		log.Fatal("WebSocket send error")
	}
	fmt.Println("Start typing your messages")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() && err == nil {
		msg := scanner.Text()
		_, err = fmt.Fprintf(ws, msg+"\n")
	}
}
