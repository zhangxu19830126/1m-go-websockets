package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"syscall"
	"time"
)

func main() {
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &syscall.Rlimit{Cur: 1000000, Max: 1000000}); err != nil {
		panic(err)
	}

	connections, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}

	u := url.URL{Scheme: "ws", Host: "172.17.0.1:8000", Path: "/"}
	log.Printf("Connecting to %s", u.String())
	var conns []*websocket.Conn
	for i := 0; i < connections; i++ {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			fmt.Println("Failed to connect", i, err)
			break
		}
		conns = append(conns, c)
		defer func() {
			c.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second))
			time.Sleep(time.Second)
			c.Close()
		}()
	}

	log.Printf("Finished initializing %d connections", len(conns))
	tts := time.Second
	if connections > 100 {
		tts = time.Millisecond * 50
	}
	for {
		time.Sleep(tts)
		idx := rand.Int() % len(conns)
		conn := conns[idx]
		log.Printf("Conn %d sending message", idx)
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Hello from conn %v", idx)))
		if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second*5)); err != nil {
			fmt.Printf("Failed to receive pong: %v", err)
		}
	}
}

// SetUlimit sets the current process ulimit soft limit to match the hard limit - ceiling (to enable more than 1024 open connections)
// Usually the hard limit is > 4000
// In order to change the hard limit, the user needs root privileges or have capability of SYS_RESOURCE
// This ulimit configuration is currently not supported by kubernetes
// https://github.com/kubernetes/kubernetes/issues/3595
