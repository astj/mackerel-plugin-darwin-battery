package main

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraphDefinition(t *testing.T) {
	var darwinBattery DarwinBatteryPlugin

	graphdef := darwinBattery.GraphDefinition()
	if len(graphdef) != 1 {
		t.Errorf("length of graphdef: %d should be 1", len(graphdef))
	}
}

func TestParse(t *testing.T) {
	var darwinBattery DarwinBatteryPlugin
	stub := `06/27/16 14:28:00
 No AC; Not Charging; 98%; Cap=6206: FCC=6286; Design=6330; Time=0:00; 0mA; Cycles=94/1000; Location=0;
 Polled boot=06/23/16 18:05:51; Full=06/27/16 14:27:17; User visible=06/27/16 14:27:17
`
	batteryStatsBuffer := bytes.NewBufferString(stub)
	stat, err := darwinBattery.ParsePmsetStats(batteryStatsBuffer)
	fmt.Println(stat)
	assert.Nil(t, err)

	// Stats
	assert.EqualValues(t, reflect.TypeOf(stat["cap"]).String(), "string")
	assert.EqualValues(t, reflect.TypeOf(stat["fcc"]).String(), "string")
	assert.EqualValues(t, reflect.TypeOf(stat["design"]).String(), "string")

	assert.EqualValues(t, stat["cap"].(string), "6206")
	assert.EqualValues(t, stat["fcc"].(string), "6286")
	assert.EqualValues(t, stat["design"].(string), "6330")
}
