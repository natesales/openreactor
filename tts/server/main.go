package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	listenAddr = flag.String("l", ":8084", "server listen address")
	piperBin   = flag.String("p", "./tts/piper/piper", "path to piper binary")
	voiceModel = flag.String("m", "./tts/voice.onnx", "path to voice model")
	verbose    = flag.Bool("v", false, "verbose logging")
)

func internalServerError(w http.ResponseWriter, err error, args ...string) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(fmt.Sprintf(
		"%s %v",
		strings.Join(args, " "), err,
	)))
}

func main() {
	flag.Parse()
	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	http.HandleFunc("/tts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Connection", "Keep-Alive")
		w.Header().Set("Transfer-Encoding", "chunked")

		for {
			cmd := exec.Command(*piperBin, "--model", *voiceModel, "--output-file", "-")
			cmd.Stdin = strings.NewReader(r.URL.Query().Get("text"))

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				internalServerError(w, err, "getting stdout pipe")
				return
			}
			if err := cmd.Start(); err != nil {
				internalServerError(w, err, "starting command")
				return
			}

			_, err = io.Copy(w, stdout)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error copying command output: %v", err), http.StatusInternalServerError)
				return
			}

			if err := cmd.Wait(); err != nil {
				internalServerError(w, err, "waiting for command")
				return
			}

			return
		}
	})

	log.WithFields(log.Fields{
		"port":  *listenAddr,
		"piper": *piperBin,
		"model": *voiceModel,
	}).Info("Starting server")
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
