package main

import (
	"encoding/json"
	"flag"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/middleware/cors"
	log "github.com/sirupsen/logrus"

	"github.com/natesales/openreactor/pkg/service"
)

var (
	listenAddr = flag.String("l", ":80", "HTTP listen address")
	verbose    = flag.Bool("v", false, "verbose logging")
)

var conns = map[*websocket.Conn]bool{}

func emit(v any) {
	log.Debugf("Emitting message: %v", v)
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

	app := service.NewApp()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // TODO: remove
	}))

	registerAlertHandlers(app.Group("/alert"))
	registerStateHandlers(app.Group("/fsm"))
	registerAPIHandlers(app.Group("/api"))

	log.WithFields(log.Fields{
		"listenAddr": *listenAddr,
	}).Info("Starting server")
	log.Fatal(app.Listen(*listenAddr))
}
