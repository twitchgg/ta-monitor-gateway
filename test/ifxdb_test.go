package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func TestInfluxDB(t *testing.T) {
	userName := "admin"
	password := "1234qwer"
	client := influxdb2.NewClient("http://10.25.133.27:8086",
		fmt.Sprintf("%s:%s", userName, password))
	wi := client.WriteAPIBlocking("", "TA-SNMP")
	p := influxdb2.NewPoint("TA-SNMP",
		map[string]string{"unit": "temperature"},
		map[string]interface{}{"avg": 24.5, "max": 45},
		time.Now())
	fmt.Println(wi.WritePoint(context.Background(), p))
}
