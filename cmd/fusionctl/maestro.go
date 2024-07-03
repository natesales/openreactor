package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type response struct {
	Msg  string          `json:"msg,omitempty"`
	Data json.RawMessage `json:"data,omitempty"`
}

func (r *response) Display() {
	fmt.Println(r.Msg)
	if r.Data == nil {
		fmt.Printf("%v", r.Data)
	}
}

func parseResponse(resp *http.Response) (*response, error) {
	var r response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	log.Debugf("Response: %+v", r)
	return &r, nil
}

func post(route string, body io.Reader) (*response, error) {
	resp, err := http.Post(
		maestroServer+"/api/"+route,
		"application/json",
		body,
	)
	if err != nil {
		return nil, err
	}

	return parseResponse(resp)
}

func get(route string) (*response, error) {
	resp, err := http.Get(maestroServer + "/api/" + route)
	if err != nil {
		return nil, err
	}

	return parseResponse(resp)
}
