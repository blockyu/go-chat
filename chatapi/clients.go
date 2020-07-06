package chatapi

import (
	"bufio"
	"io"
	"log"
)

type client struct {
	* bufio.Reader
	* bufio.Writer
	wc chan string
}

func StartClient(name string, msgCh chan<- string, cn io.ReadWriteCloser, roomName string) (chan<- string, <-chan struct{}) {
	c := new(client)
	c.Reader = bufio.NewReader(cn)
	c.Writer = bufio.NewWriter(cn)
	c.wc = make(chan string)
	channelDone := make(chan struct{})

	go func() {
		scanner := bufio.NewScanner(c.Reader)
		for scanner.Scan(){
			msg := name + ":" + scanner.Text() + "\n"
			log.Printf("%s|%s", roomName, msg)
			msgCh <- msg
		}
		close(channelDone)
		cn.Close()
	}()
}

func (c *client) writerMonitor() {
	go func() {
		for s := range c.wc {
			c.WriteString(s)
			c.Flush()
		}
	}()
}
