package main

import (
	"encoding/json"
	"flag"
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var (
	listenAddr = flag.String("l", ":80", "HTTP listen address")
	verbose    = flag.Bool("v", false, "verbose logging")
)

var conns = map[*websocket.Conn]bool{}

func emit(v any) {
	j, err := json.Marshal(v)
	if err != nil {
		log.Errorf("Error marshalling message: %v", err)
		return
	}
	for client := range conns {
		if err := client.WriteMessage(websocket.TextMessage, j); err != nil {
			log.Debugf("Error writing message to client: %v", err)
			_ = client.Close()
			delete(conns, client)
		}
	}
}

func main() {
	flag.Parse()
	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	mux := http.NewServeMux()

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer c.Close()

		conns[c] = true

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Debugf("reading from client: %s", err)
				break
			}
			log.Infof("recv: %s", message)
		}
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("alert server ok"))
	})

	mux.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		emit(map[string]string{
			"type":    "logMessage",
			"message": r.URL.Query().Get("msg"),
		})
	})

	mux.HandleFunc("/alert", func(w http.ResponseWriter, r *http.Request) {
		emit(map[string]string{
			"type": "audioAlert",
			"text": r.URL.Query().Get("msg"),
		})
	})

	log.WithFields(log.Fields{
		"listenAddr": *listenAddr,
	}).Info("Starting server")
	log.Fatal(http.ListenAndServe(*listenAddr, mux))
}
