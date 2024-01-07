package main

import (
	"flag"
	"log"

	"github.com/natesales/openreactor/config"
)

var (
	configFile = flag.String("c", "", "config file")
)

func main() {
	flag.Parse()

	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatal(err)
	}
}
