package service

import (
	"flag"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/natesales/openreactor/pkg/serial"
)

type pollFunc func() error

type Service struct {
	SerialPort   *serial.Port
	pollFunc     pollFunc
	pollInterval time.Duration
	app          *fiber.App
	listenAddr   string
	Log          *logrus.Entry
}

func New(baud int) *Service {
	var (
		serialPort   = flag.String("s", "", "Serial port")
		listenAddr   = flag.String("l", ":80", "API listen address")
		pollInterval = flag.Duration("i", 1*time.Second, "Poll interval")
		verbose      = flag.Bool("v", false, "Enable verbose logging")
		trace        = flag.Bool("vv", false, "Enable trace logging")
	)

	flag.Parse()
	if *verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if *trace {
		logrus.SetLevel(logrus.TraceLevel)
	}

	// Connect to serial port
	p := serial.New(*serialPort, baud)
	if err := p.Connect(); err != nil {
		logrus.Fatalf("serial connect: %s", err)
	}

	return &Service{
		SerialPort:   p,
		pollFunc:     nil,
		pollInterval: *pollInterval,
		app: fiber.New(fiber.Config{
			DisableStartupMessage: true,
		}),
		listenAddr: *listenAddr,
		Log:        logrus.WithField("svc", os.Args[0]), // TODO: Maybe get service name elsewhere
	}
}

// Start starts the metrics poller and API server
func (s *Service) Start() {
	ticker := time.NewTicker(s.pollInterval)
	go func() {
		for ; true; <-ticker.C {
			if err := s.pollFunc(); err != nil {
				s.Log.Warnf("polling: %s", err)
			}
		}
	}()

	if err := s.app.Listen(s.listenAddr); err != nil {
		s.Log.Fatalf("app listen: %s", err)
	}
}

func (s *Service) SetPoller(p pollFunc) {
	s.pollFunc = p
}
