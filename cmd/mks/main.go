package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/natesales/openreactor/pkg/db"
	"github.com/natesales/openreactor/pkg/service"
)

var svc = service.New(9600)

func exec(cmd string) (string, error) {
	_, err := svc.SerialPort.P.Write([]byte(fmt.Sprintf("@254%s;FF", cmd)))
	if err != nil {
		svc.Log.Fatal(err)
	}

	buf := make([]byte, 0)

	for {
		b := make([]byte, 1)
		_, err := svc.SerialPort.P.Read(b)
		if err != nil {
			svc.Log.Warnf("reading from serial port: %v", err)
			continue
		}
		buf = append(buf, b[0])

		if buf[len(buf)-1] == 'F' && buf[len(buf)-2] == 'F' {
			break
		}
	}

	out := string(buf)
	out = strings.TrimLeft(out, "@253")
	out = strings.TrimRight(out, ";FF")

	return out, nil
}

func readFloat(cmd string) (float64, error) {
	resp, err := exec("PR1")
	if err != nil {
		return -1, err
	}
	if strings.HasPrefix(resp, "ACK") {
		resp = strings.TrimLeft(resp, "ACK")
		return strconv.ParseFloat(resp, 64)
	} else {
		return -1, fmt.Errorf("NAK: %s", resp)
	}
}

func main() {
	if err := svc.SerialPort.Flush(); err != nil {
		svc.Log.Fatal(err)
	}

	svc.SetPoller(func() error {
		pirani, err := readFloat("PR1")
		if err != nil {
			return err
		}
		piezo, err := readFloat("PR2")
		if err != nil {
			return err
		}
		if err := db.Write("mks", nil, map[string]any{
			"pirani": pirani,
			"piezo":  piezo,
		}); err != nil {
			return err
		}

		return nil
	})

	svc.Start()
}
