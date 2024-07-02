package db

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/query"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	log "github.com/sirupsen/logrus"
)

var (
	client   influxdb2.Client
	writeAPI api.WriteAPIBlocking
	queryAPI api.QueryAPI
)

const (
	org    = "reactor"
	bucket = "dlog"
)

type Measurement string

var (
	HVSetpoint Measurement = "hv_setpoint"
	HVCurrent  Measurement = "hv_current"
	HVVoltage  Measurement = "hv_voltage"

	TurboRunning Measurement = "turbo_running"
	TurboSpeed   Measurement = "turbo_hz"
	TurboCurrent Measurement = "turbo_current"

	EdwardsGaugeTorr Measurement = "edwards_torr"
	EdwardsGaugeVolt Measurement = "edwards_volt"

	MKSGaugeVacuum Measurement = "mksgauge"

	NeutronCPS Measurement = "neutron_cps"

	MKSMFCFlow     Measurement = "mksmfc_flow"
	MKSMFCSetPoint Measurement = "mksmfc_setpoint"

	SierraMFCFlow     Measurement = "sierramfc_flow"
	SierraMFCSetPoint Measurement = "sierramfc_setpoint"
)

func init() {
	c := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	options := influxdb2.DefaultOptions()
	options.SetHTTPClient(c)
	client = influxdb2.NewClientWithOptions(
		os.Getenv("OPENREACTOR_INFLUXDB_URL"),
		os.Getenv("OPENREACTOR_INFLUXDB_TOKEN"),
		options,
	)
	writeAPI = client.WriteAPIBlocking(org, bucket)
	queryAPI = client.QueryAPI(org)
}

// Write a measurement to the database
func Write(measurement Measurement, tags map[string]string, fields map[string]any) error {
	point := write.NewPoint(string(measurement), tags, fields, time.Now())
	return writeAPI.WritePoint(context.Background(), point)
}

// Last gets the latest point from the database
func Last(measurement Measurement) (*query.FluxRecord, error) {
	q := `from(bucket: "` + bucket + `")
	|> range(start: -7d)
	|> filter(fn: (r) => r._measurement == "` + string(measurement) + `")
	|> last()`

	log.Debugf("Executing query: %s", q)
	result, err := queryAPI.Query(context.Background(), q)
	if err != nil {
		return nil, err
	}

	for result.Next() {
		return result.Record(), nil
	}

	return nil, fmt.Errorf("no results found for %s", measurement)
}

// LastOrNil gets the latest point from the database, or nil if there is an error
func LastOrNil(measurement Measurement) any {
	point, err := Last(measurement)
	if err != nil {
		log.Debugf("Error getting last point for %s: %s", measurement, err)
		return nil
	}

	return point.Value()
}

// LastFloat gets the latest point from the database as a float64
func LastFloat(measurement Measurement) (float64, error) {
	point, err := Last(measurement)
	if err != nil {
		return -1, err
	}

	return point.Value().(float64), nil
}
