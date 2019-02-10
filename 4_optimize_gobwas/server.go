package main

import (
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"syscall"
)

var epoller *epoll

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade connection
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return
	}
	if err := epoller.Add(conn); err != nil {
		log.Printf("Failed to add connection")
		conn.Close()
	}
	//// Read messages from socket
	//for {
	//	_, msg, err := conn.ReadMessage()
	//	if err != nil {
	//		log.Printf("Failed to read message %v", err)
	//		return
	//	}
	//	log.Println(string(msg))
	//}
}

func main() {
	// Increase resources limitations
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &syscall.Rlimit{Cur: 1000000, Max: 1000000}); err != nil {
		panic(err)
	}

	// Enable pprof hooks
	go func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("Pprof failed:", err)
		}
	}()

	// Start epoll
	var err error
	epoller, err = MkEpoll()
	if err != nil {
		panic(err)
	}

	go Start()

	http.HandleFunc("/", wsHandler)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}

func Start() {
	for {
		connections, err := epoller.Wait()
		if err != nil {
			log.Printf("Failed to epoll wait %v", err)
			continue
		}
		for _, conn := range connections {
			if conn == nil {
				break
			}
			_, _, err := wsutil.ReadClientData(conn)
			//msgs, err := wsutil.ReadMessage(conn, ws.StateServerSide, nil)
			if err != nil {
				// handle error
			}

			if err != nil {
				log.Printf("Failed to read message %v", err)
				if err := epoller.Remove(conn); err != nil {
					log.Printf("Failed to remove %v", err)
				}
			} else {
				//log.Printf("msg: %s", string(msg))
			}
		}
	}
}
