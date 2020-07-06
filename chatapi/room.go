package chatapi

import (
	"io"
	"log"
	"sync"
)

type Room struct {
	name 	string
	clients map[string]chan<- string
	Msgch chan string
	Quit chan struct{}
	*sync.RWMutex
}

func CreateRoom(rname string) *Room {
	r := &Room{
		name: rname,
		Msgch: make(chan string),
		RWMutex: new(sync.RWMutex),
		clients: make(map[string]chan<- string),
		Quit:	make(chan struct{}),
	}
	r.Run()
	return r
}

func (r *Room) AddClient(c io.ReadWriteCloser, clientName string) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.clients[clientName]; ok {
		log.Printf("Client %s already exist in chat room %s, existing...", clientName, r.name)
		return
	}
	log.Printf("Adding client %s \n", clientName)
	wc, done := StartClient(clientName, r.Msgch, c, r.name)
	r.clients[clientName] = wc

	// remove client when done is signalled
	go func() {
		<- done
		r.RemoveClientSync(clientName)
	}()
}

func (r *Room) ClCount() int {
	return len(r.clients)
}

func (r *Room) RemoveClientSync(name string) {
	r.Lock()
}

func (r *Room) Run() {
	log.Println("Starting chat room", r.name)
	// handle the chat room, main message channel
	go func() {
		for msg := range r.Msgch {
			r.broadcastMsg(msg)
		}
	}()

	// handle when the quit channel is triggered
	go func() {
		<- r.Quit
		r.CloseChatRoomSync()
	}()

}

//CloseChatRoomSync closes a chat room. This is a blocking call
func (r *Room) CloseChatRoomSync() {
	r.Lock()
	defer r.Unlock()
	close(r.Msgch)
	for name := range r.clients {
		delete(r.clients, name)
	}
}

//fan out is used to distribute the chat message
func (r *Room) broadcastMsg(msg string) {
	r.RLock()
	defer r.RUnlock()
	for _, wc := range r.clients {
		go func(wc chan<- string) {
			wc <- msg
		}(wc)
	}
}
