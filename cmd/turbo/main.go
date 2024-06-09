package main

import (
	"net/http"

	"github.com/natesales/openreactor/pkg/alert"
	"github.com/natesales/openreactor/pkg/db"
	"github.com/natesales/openreactor/pkg/service"
)

func main() {
	t := TCP015{}
	svc := service.New(t)
	if err := t.Connect(); err != nil {
		svc.Log.Fatal(err)
	}

	fw, err := t.FirmwareVersion()
	if err != nil {
		svc.Log.Fatal(err)
	}
	svc.Log.Infof("Turbo pump %s", fw)

	svc.SetPoller(func() error {
		hz, err := t.Hz()
		if err != nil {
			return err
		}

		current, err := t.CurrentDraw()
		if err != nil {
			return err
		}

		if err := db.Write("turbo_hz", nil, map[string]any{"hz": hz}); err != nil {
			return err
		}
		if err := db.Write("turbo_current", nil, map[string]any{"current": current}); err != nil {
			return err
		}

		isRunning, err := t.IsRunning()
		if err != nil {
			return err
		}
		isRunningInt := 0
		if isRunning {
			isRunningInt = 1
		}

		if err := db.Write("turbo_running", nil, map[string]any{"running": isRunningInt}); err != nil {
			return err
		}
	})

	http.HandleFunc("/on", func(w http.ResponseWriter, r *http.Request) {
		alert.Alert("Starting turbo")
		if err := t.On(); err != nil {
			w.Write([]byte("Error: " + err.Error()))
		}
		w.Write([]byte("ok"))
	})
	http.HandleFunc("/off", func(w http.ResponseWriter, r *http.Request) {
		alert.Log("Stopping turbo")
		if err := t.Off(); err != nil {
			w.Write([]byte("Error: " + err.Error()))
		}
		w.Write([]byte("ok"))
	})

	svc.Start()
}
