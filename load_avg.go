package monitor

import (
	"io/ioutil"
	"strconv"
	"strings"
)

// LoadAvg returns the first 3 values in /proc/loadavg, as floats
type LoadAvg struct {
	OneMinute, FiveMinutes, FifteenMinutes float64
}

func newLoadAvg(avgs [3]float64) *LoadAvg {
	return &LoadAvg{
		OneMinute:      avgs[0],
		FiveMinutes:    avgs[1],
		FifteenMinutes: avgs[2],
	}
}

/*
Source: https://github.com/nodequery/nq-agent/blob/master/nq-agent.sh
	/proc/cpuinfo
	/proc/meminfo
	/proc/stat
	/proc/uptime
	/proc/loadavg
*/

/*
Get memory stats

$ cat /proc/meminfo | grep ^Mem

MemTotal:        4052056 kB
MemFree:          585140 kB
MemAvailable:    1712788 kB
*/

// GetLoadAvg returns the current load avg
func GetLoadAvg() (*LoadAvg, error) {
	loadavg, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		return &LoadAvg{}, err
	}
	strAvgs := strings.Split(string(loadavg), " ")

	avgs, err := convertToFloats(strAvgs)
	if err != nil {
		return &LoadAvg{}, err
	}

	return newLoadAvg(avgs), err
}

func convertToFloats(strAvgs []string) ([3]float64, error) {
	var (
		f      float64
		err    error
		floats [3]float64
	)

	if len(strAvgs) != 5 {
		panic("convertToFloats: invalid loadavg!")
	}

	for i := 0; i < len(floats); i++ {
		f, err = strconv.ParseFloat(strAvgs[i], 64)
		if err != nil {
			return floats, err
		}
		floats[i] = f
	}
	return floats, err
}
