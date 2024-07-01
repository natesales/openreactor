package service

import (
	"fmt"
	"net/http"
)

// RPC sends a request to a remote service
func RPC(route string) error {
	resp, err := http.Get("http://" + route)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
