package main

import (
	"log"

	"github.com/benjamin-thomas/monitor"
)

var (
	maxAvg = 0.4
)

var colors = struct {
	none, red, green string
}{
	none:  "\033[1;m",
	red:   "\033[1;31m",
	green: "\033[1;32m",
}

func colorize(str, colorCode string) string {
	return colorCode + str + colors.none
}

func main() {
	avg, err := monitor.GetLoadAvg()
	if err != nil {
		log.Fatal(err)
	}
	if avg.OneMinute < maxAvg {
		log.Printf(colorize("Avg is OK (curr=%.2f, maxAvg=%.2f)", colors.green), avg.OneMinute, maxAvg)
	} else {
		log.Printf(colorize("Avg is too high!! (curr=%.2f, maxAvg=%.2f)", colors.red), avg.OneMinute, maxAvg)
	}
	log.Printf("%#v\n", avg)
}
