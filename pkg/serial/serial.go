package serial

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

func New(port string, baud int) *Port {
	return &Port{
		Port: port,
		Baud: baud,
		Lock: sync.Mutex{},
	}
}

type Port struct {
	Port string
	Baud int

	P    serial.Port
	Lock sync.Mutex
}

// Connect connects to the serial port
func (p *Port) Connect() error {
	mode := &serial.Mode{
		BaudRate: p.Baud,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}
	var err error
	log.Infof("Connecting to %s", p.Port)
	p.P, err = serial.Open(p.Port, mode)
	return err
}

// Close closes the serial port
func (p *Port) Close() error {
	return p.P.Close()
}

// Reconnect closes and reopens the serial port
func (p *Port) Reconnect() error {
	if err := p.Close(); err != nil {
		return err
	}
	return p.Connect()
}

// Flush purges the input and output buffers
func (p *Port) Flush() error {
	if err := p.P.ResetInputBuffer(); err != nil {
		return err
	}
	if err := p.P.ResetOutputBuffer(); err != nil {
		return err
	}
	return nil
}

func (p *Port) Write(message []byte) error {
	p.Lock.Lock()
	_, err := p.P.Write(message)
	p.Lock.Unlock()
	return err
}

func (p *Port) Read(b []byte) (n int, err error) {
	p.Lock.Lock()
	defer p.Lock.Unlock()
	return p.P.Read(b)
}

// Send sends a message and returns the response
func (p *Port) Send(message []byte) (string, error) {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	if err := p.Flush(); err != nil {
		return "", err
	}

	_, err := p.P.Write(message)
	if err != nil {
		return "", err
	}

	buf := make([]byte, 0)

	for {
		b := make([]byte, 1)
		_, err := p.P.Read(b)
		if err != nil {
			log.Warnf("reading from serial port: %v", err)
			continue
		}

		if b[0] == '\r' {
			break
		}

		buf = append(buf, b[0])
	}

	return string(buf), nil
}
