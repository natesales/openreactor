package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Build struct {
	Context string
	Args    []string
}

type Container struct {
	Image   string
	Build   Build
	Devices []string
	Ports   []string
}

func main() {
	svcFilename := flag.String("f", "svc.yml", "Service file")
	flag.Parse()

	var svc struct {
		Ports map[string]string `yaml:"ports"`
	}
	svcFile, err := os.ReadFile(*svcFilename)
	if err != nil {
		log.Fatal(err)
	}
	if err := yaml.Unmarshal(svcFile, &svc); err != nil {
		log.Fatal(err)
	}

	outFile, err := os.Create("docker-compose.svc.yml")
	if err != nil {
		log.Fatal(err)
	}

	var cfg struct {
		Services map[string]Container `yaml:"services"`
	}
	cfg.Services = make(map[string]Container)

	i := 0
	for name, port := range svc.Ports {
		cfg.Services[name] = Container{
			Image: "openreactor-" + name,
			Build: Build{
				Context: ".",
				Args:    []string{"SVC=" + name},
			},
			Devices: []string{port + ":/serial"},
			Ports:   []string{fmt.Sprintf("%d:80", 8080+i)},
		}
		log.Infof("Added service %s", name)
		i++
	}

	if err := yaml.NewEncoder(outFile).Encode(cfg); err != nil {
		log.Fatal(err)
	}
}
