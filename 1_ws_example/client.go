package main

import (
	"fmt"
	"github.com/gorilla/websocket"
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

	u := url.URL{Scheme: "ws", Host: "172.17.0.1:8000", Path: "/ws"}
	fmt.Println("connecting to", u.String())

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

	fmt.Println("Finished initializing connections", len(conns))
	for i := 0; i < len(conns); i++ {
		time.Sleep(time.Millisecond * 1000)
		idx := rand.Int() % len(conns)
		conn := conns[idx]
		fmt.Println("Conn sending message", idx)
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Hello from conn %v", idx)))
		if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second * 5)); err != nil {
			fmt.Println("Error receiving pong", err)
		}
	}
	time.Sleep(time.Minute * 3)
}


// SetUlimit sets the current process ulimit soft limit to match the hard limit - ceiling (to enable more than 1024 open connections)
// Usually the hard limit is > 4000
// In order to change the hard limit, the user needs root privileges or have capability of SYS_RESOURCE
// This ulimit configuration is currently not supported by kubernetes
// https://github.com/kubernetes/kubernetes/issues/3595
