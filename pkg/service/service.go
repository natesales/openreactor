package service

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"

	"github.com/natesales/openreactor/pkg/serial"
)

type pollFunc func() error

type Service struct {
	SerialPort *serial.Port
	App        *fiber.App
	Log        *logrus.Entry

	listenAddr   string
	pollFunc     pollFunc
	pollInterval time.Duration
}

func New(baud int) *Service {
	var (
		serialPort   = flag.String("s", "/serial", "Serial port")
		listenAddr   = flag.String("l", ":80", "API listen address")
		pollInterval = flag.Duration("i", 1*time.Second, "Poll interval")
		verbose      = flag.Bool("v", false, "Enable verbose logging")
		trace        = flag.Bool("vv", false, "Enable trace logging")
	)

	// Parse name from binary name
	nameParts := strings.Split(os.Args[0], "/")
	s := nameParts[len(nameParts)-1]
	logger := logrus.WithField("svc", s)

	flag.Parse()
	if *verbose {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if *trace {
		logrus.SetLevel(logrus.TraceLevel)
	}

	if *serialPort == "" {
		logger.Fatalf("required flag -s not provided")
	}

	// Connect to serial port
	p := serial.New(*serialPort, baud)
	logger.Infof("Connecting to %s", *serialPort)
	if err := p.Connect(); err != nil {
		logger.Fatalf("serial connect: %s", err)
	}

	return &Service{
		SerialPort: p,
		Log:        logger,
		App: fiber.New(fiber.Config{
			DisableStartupMessage: true,
		}),
		pollFunc:     nil,
		pollInterval: *pollInterval,
		listenAddr:   *listenAddr,
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

	if err := s.App.Listen(s.listenAddr); err != nil {
		s.Log.Fatalf("app listen: %s", err)
	}
}

func (s *Service) SetPoller(p pollFunc) {
	s.pollFunc = p
}
