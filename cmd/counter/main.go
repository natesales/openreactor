package main

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/natesales/openreactor/pkg/db"
	"github.com/natesales/openreactor/pkg/service"
)

func main() {
	svc := service.New(115200)

	last := 0
	svc.App.Get("/last", func(ctx *fiber.Ctx) error {
		return ctx.SendString(strconv.Itoa(last))
	})

	go svc.Start()

	buf := make([]byte, 0)
	for {
		b := make([]byte, 1)
		_, err := svc.SerialPort.Read(b)
		if err != nil {
			svc.Log.Warnf("reading from serial port: %v", err)
			continue
		}

		if b[0] == ';' {
			line := strings.TrimSpace(string(buf))
			count, err := strconv.Atoi(line)
			if err != nil {
				svc.Log.Warnf("parsing float: %v", err)
			}
			svc.Log.Debugf("%d cps", count)

			last = count
			if err := db.Write(db.NeutronCPS, nil, map[string]any{"cps": count}); err != nil {
				svc.Log.Warn(err)
				continue
			}

			buf = make([]byte, 0)
		} else {
			buf = append(buf, b[0])
		}
	}
}
