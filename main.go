package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	socketio "github.com/googollee/go-socket.io"
)

var (
	mu sync.Mutex
)

func main() {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "message", func(s socketio.Conn, room string) {
		s.Join(room)
		go func() {
			for {
				mu.Lock()
				randomNumber := rand.Intn(100) + 1
				println("random", randomNumber)
				server.BroadcastToRoom("", room, "broad", randomNumber)
				time.Sleep(1 * time.Second)
				mu.Unlock()
			}
		}()

	})

	server.OnError("/", func(s socketio.Conn, e error) {
		//server.Remove(s.ID())
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		// Add the Remove session id. Fixed the connection & mem leak
		//server.Remove(s.ID())
		fmt.Println("closed", reason)
	})

	go server.Serve()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	log.Println("Serving at localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
