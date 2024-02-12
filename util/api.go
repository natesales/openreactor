package util

import "net/http"

func HandleExec(f func() error) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(); err != nil {
			w.Write([]byte("Error: " + err.Error()))
		}
		w.Write([]byte("ok"))
	}
}
