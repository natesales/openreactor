package ws

import (
	"encoding/json"
	"fmt"
)

type wsHandler func(msg string) error

var callbacks = map[string]wsHandler{}

func HandleFunc(name string, callback wsHandler) {
	callbacks[name] = callback
}

func Handle(msg []byte) error {
	var m struct {
		Name    string `json:"name"`
		Payload string `json:"payload"`
	}
	if err := json.Unmarshal(msg, &m); err != nil {
		return err
	}

	cb, ok := callbacks[m.Name]
	if !ok {
		return fmt.Errorf("no handler for message: %s", m.Name)
	}
	return cb(m.Payload)
}
