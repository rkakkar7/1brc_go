package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
)

type data struct {
	minTemp float64
	maxTemp float64
	sum     float64
	count   int64
}

func processData(file io.Reader, outputWriter io.Writer) {
	fileScanner := bufio.NewScanner(file)
	stationData := make(map[string]*data)
	// only reading 1 line per scan
	for fileScanner.Scan() {
		line := fileScanner.Text()
		splitString := strings.Split(line, ";")

		if len(splitString) < 2 {
			continue
		}
		station := splitString[0]
		temp, err := strconv.ParseFloat(splitString[1], 64)
		if err != nil {
			log.Fatalf("error converting string to float: %+v", err)
		}

		if val, ok := stationData[station]; ok {
			val.maxTemp = max(val.maxTemp, temp)
			val.minTemp = min(val.minTemp, temp)
			val.sum += temp
			val.count++
		} else {
			val = &data{}
			val.count = 1
			val.maxTemp = temp
			val.minTemp = temp
			val.sum = temp
		}
	}
	stationSlice := make([]string, 0)
	for station := range stationData {
		stationSlice = append(stationSlice, station)
	}
	sort.Slice(stationSlice, func(i, j int) bool {
		return stationSlice[i] < stationSlice[j]
	})

	fmt.Fprintf(outputWriter, "{")
	for i, station := range stationSlice {
		if i > 0 {
			fmt.Fprintf(outputWriter, ", ")
		}
		s := stationData[station]
		avg := round(round(s.sum) / float64(s.count))
		fmt.Fprintf(outputWriter, "%s=%.1f/%.1f/%.1f", station, stationData[station].minTemp, avg, stationData[station].maxTemp)
	}
	fmt.Fprintf(outputWriter, "}\n")
}

// rounding floats to 1 decimal place with 0.05 rounding up to 0.1
func round(x float64) float64 {
	return math.Floor((x+0.05)*10) / 10
}

// FLAGS
// -file=file location
// -cpu turn on cpu profiling

func main() {
	fileLocation := flag.String("file", "", "file location")
	cpuProfile := flag.Bool("cpu", false, "turn on cpu profiling")
	memProfile := flag.Bool("mem", false, "turn on profiling")
	flag.Parse()

	if *cpuProfile {
		cpuProf, err := os.Create("cpuProfile")
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(cpuProf)
		defer pprof.StopCPUProfile()
	}

	if *memProfile {
		memProf, err := os.Create("cpuProfile")
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(memProf)
		defer memProf.Close()
	}

	absPath, err := filepath.Abs(*fileLocation)
	if err != nil {
		log.Fatalf("filepath: %+v", err)
	}

	file, err := os.Open(absPath)
	if err != nil {
		log.Fatalf("err: %+v", err)
	}
	defer file.Close()

	outputWriter := bufio.NewWriter(os.Stdout)
	processData(file, outputWriter)

	outputWriter.Flush()
}
