package db

import (
	"context"
	"time"

	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

var (
	client = influxdb2.NewClient(
		"https://influxdb-reactor.westland.as34553.net",
		"mLTGAlKl0WfX02H6C-adeG_ZeqSvzxf9LoFp3Guhj1EwgNt1mmf5UHw_dFXxWdjk90Y9a3eVYOm02aTWbD-kfg==",
	)
	writeAPI = client.WriteAPIBlocking("reactor", "dlog")
)

// Write a measurement to the database
func Write(measurement string, tags map[string]string, fields map[string]any) error {
	point := write.NewPoint(measurement, tags, fields, time.Now())
	return writeAPI.WritePoint(context.Background(), point)
}
