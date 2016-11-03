package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os/exec"
	"regexp"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

type DarwinBatteryPlugin struct {
	Prefix string
}

// GraphDefinition interface for mackerelplugin
func (d DarwinBatteryPlugin) GraphDefinition() map[string](mp.Graphs) {
	return map[string](mp.Graphs){
		"": mp.Graphs{
			Label: "Battery Capacity",
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "cap", Label: "Cap", Type: "uint64"},
				mp.Metrics{Name: "fcc", Label: "FCC", Type: "uint64"},
				mp.Metrics{Name: "design", Label: "Design", Type: "uint64"},
			},
		},
	}
}

// MetricKeyPrefix is implementation for PluginWithPrefix interface
func (d DarwinBatteryPlugin) MetricKeyPrefix() string {
	return d.Prefix
}

func (d DarwinBatteryPlugin) getStatsReaderFromPmset() (map[string]interface{}, error) {
	cmd := exec.Command("pmset", "-g", "rawbatt")
	stdout, err := cmd.StdoutPipe()
	reader := bufio.NewReader(stdout)
	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("faild to fetch pmset: %s", err)
	}
	stat, parseErr := d.ParsePmsetStats(reader)
	if parseErr != nil {
		return nil, fmt.Errorf("faild to parse pmset: %s", err)
	}
	_ = cmd.Wait()
	return stat, nil
}

func (d DarwinBatteryPlugin) ParsePmsetStats(r io.Reader) (map[string]interface{}, error) {
	stat := make(map[string]interface{})
	group := [](string){}
	scanner := bufio.NewScanner(r)
	i := 0
	for scanner.Scan() {
		i++
		// we need only 2nd line!
		if i == 2 {
			line := scanner.Text()
			s := string(line)
			assined := regexp.MustCompile(`Cap=(\d+): FCC=(\d+); Design=(\d+);`)
			group = assined.FindStringSubmatch(s)
			break
		}
	}

	stat["cap"] = group[1]
	stat["fcc"] = group[2]
	stat["design"] = group[3]
	return stat, nil
}

// FetchMetrics interface for mackerelplugin
func (d DarwinBatteryPlugin) FetchMetrics() (map[string]interface{}, error) {
	stat, err := d.getStatsReaderFromPmset()
	if err != nil {
		return nil, fmt.Errorf("Faild to fetch battery metrics: %s", err)
	}
	return stat, nil
}

func main() {
	optPrefix := flag.String("metric-key-prefix", "battery-capacity", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	d := DarwinBatteryPlugin{
		Prefix: *optPrefix,
	}
	helper := mp.NewMackerelPlugin(d)
	helper.Tempfile = *optTempfile
	helper.Run()
}
