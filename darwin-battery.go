package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"

	mp "github.com/mackerelio/go-mackerel-plugin-helper"
)

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

func getStatsReaderFromPmset() (io.Reader, error) {
	cmd := exec.Command("pmset", "-g", "rawbatt")
	stdout, err := cmd.StdoutPipe()
	reader := bufio.NewReader(stdout)
	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("faild to fetch pmset: %s", err)
	}
	return reader, nil
}

func ParsePmsetStats(r io.Reader) (map[string]interface{}, error) {
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
		}
	}

	stat["cap"] = group[1]
	stat["fcc"] = group[2]
	stat["design"] = group[3]
	return stat, nil
}

// FetchMetrics interface for mackerelplugin
func (d DarwinBatteryPlugin) FetchMetrics() (map[string]interface{}, error) {
	reader, err := getStatsReaderFromPmset()
	if err != nil {
		return nil, fmt.Errorf("Faild to fetch battery metrics: %s", err)
	}
	stat, err := ParsePmsetStats(reader)
	if err != nil {
		return nil, fmt.Errorf("Faild to parse battery metrics: %s", err)
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
	if helper.Tempfile == "" {
		helper.Tempfile = fmt.Sprintf("/tmp/mackerel-plugin-%s", *optPrefix)
	}
	helper.Run()
}
