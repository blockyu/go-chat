package websocket

import (
	"fmt"
	"github.com/blockyu/go-chat/chatapi"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type websocketConnWrapper struct {
	*websocket.Conn
	msgType int
}

func (wsWrapper *websocketConnWrapper) Read(b []byte) (int, error) {
	t, r, err := wsWrapper.NextReader()
	wsWrapper.msgType = t
	if err != nil {
		log.Println("Websocket err:", err)
		return 0, err
	}
	return r.Read(b)
}

func (wsWrapper *websocketConnWrapper) Write(p []byte) (int, error) {
	err := wsWrapper.WriteMessage(wsWrapper.msgType, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (wsWrapper *websocketConnWrapper) Close() error {
	return wsWrapper.Conn.Close()
}

type wsChatHandler struct {
	*websocket.Upgrader
	*chatapi.ChatAPI
}

func RunWSChatWithExistingAPI(url, address string, Rbsize, Wbsize int, chat *chatapi.ChatAPI) error {
	wsHandler := &wsChatHandler{
		Upgrader: &websocket.Upgrader{
			ReadBufferSize: Rbsize,
			WriteBufferSize: Wbsize,
		},
		ChatAPI: chat,
	}
	http.HandleFunc(url, wsHandler.wshandler)
	return http.ListenAndServe(address, nil)
}

func RunWSChat(url, address string, Rbsize, Wbsize int) error {
	return RunWSChatWithExistingAPI(url, address, Rbsize, Wbsize, chatapi.NewChat())
}

func (wh *wsChatHandler) wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wh.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error occured while trying to upgrade websocket", err)
		fmt.Fprintf(w, "Error occured while trying to upgrade websocket: %v", err)
		return
	}
	wsWrapper := &websocketConnWrapper{Conn: conn}
	wh.AddClient(wsWrapper)
}