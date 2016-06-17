package main

import (
	"bufio"
	"flag"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

// UptimePlugin mackerel plugin
type DarwinBatteryPlugin struct {
	Prefix string
}

// GraphDefinition interface for mackerelplugin
func (d DarwinBatteryPlugin) GraphDefinition() map[string](mp.Graphs) {
	labelPrefix := strings.Title(d.Prefix)

	return map[string](mp.Graphs){
		d.Prefix: mp.Graphs{
			Label: labelPrefix,
			Unit:  "integer",
			Metrics: [](mp.Metrics){
				mp.Metrics{Name: "cap", Label: "Cap", Type: "uint64"},
				mp.Metrics{Name: "fcc", Label: "FCC", Type: "uint64"},
				mp.Metrics{Name: "design", Label: "Design", Type: "uint64"},
			},
		},
	}
}

func getMetricsFromPmset() (string, string, string, error) {
	cmd := exec.Command("pmset", "-g", "rawbatt")
	stdout, err := cmd.StdoutPipe()
	scanner := bufio.NewScanner(stdout)
	err = cmd.Start()
	if err != nil {
		return "0", "0", "0", fmt.Errorf("faild to fetch pmset: %s", err)
	}
	i := 0
	group := [](string){}
	for scanner.Scan() {
		i++
		// we need only 2nd line!
		if i == 2 {
			line := scanner.Text()
			s := string(line)
			assined := regexp.MustCompile(`Cap=(\d+): FCC=(\d+); Design=(\d+);`)
			group = assined.FindStringSubmatch(s)

		}
	}
	return group[1], group[2], group[3], nil
}

// FetchMetrics interface for mackerelplugin
func (d DarwinBatteryPlugin) FetchMetrics() (map[string]interface{}, error) {
	stat := make(map[string]interface{})
	cap, fcc, design, err := getMetricsFromPmset()
	if err != nil {
		return nil, fmt.Errorf("Faild to fetch uptime metrics: %s", err)
	}
	stat["cap"] = cap
	stat["fcc"] = fcc
	stat["design"] = design
	return stat, nil
}

func main() {
	optPrefix := flag.String("metric-key-prefix", "battery capacity", "Metric key prefix")
	optTempfile := flag.String("tempfile", "", "Temp file name")
	flag.Parse()

	d := DarwinBatteryPlugin{
		Prefix: *optPrefix,
	}
	helper := mp.NewMackerelPlugin(d)
	helper.Tempfile = *optTempfile
	if helper.Tempfile == "" {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-%s", *optPrefix)
	}
	helper.Run()
}
