package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
)

const (
	defaultNumBins = 20
	width          = 80
	height         = 24
)

var options struct {
	logX    bool
	logY    bool
	numBins int
}

func init() {
	flag.BoolVar(&options.logX, "log-x", false, "use log x-axis bins")
	flag.BoolVar(&options.logY, "log-y", false, "use log y-axis bins")
	flag.IntVar(&options.numBins, "num-bins", defaultNumBins,
		"number of bins to include (legal values: 10, 20, 40, 80)")
	flag.Parse()
}

func main() {
	binNumOk := false
	for _, validBinSize := range []int{10, 20, 40, 80} {
		if options.numBins == validBinSize {
			binNumOk = true
			break
		}
	}
	if !binNumOk {
		fmt.Printf("Illegal number of bins: %d", options.numBins)
		os.Exit(1)
	}

	_, vals, err := makeBins(options.logX, options.logY, options.numBins)
	if err != nil {
		fmt.Printf("Error reading from stdin: %v\n", err)
		os.Exit(1)
	}
	graphVals(vals)
}

func makeBins(logX, logY bool, numBins int) ([]float64, []float64, error) {
	midpoints := make([]float64, numBins)
	vals := make([]float64, numBins)
	rawData, err := readRawData()
	if err != nil {
		return midpoints, vals, err
	}
	if len(rawData) == 0 {
		return midpoints, vals, fmt.Errorf("no data parsed")
	}

	// figure out range
	var minVal, maxVal float64
	for i, val := range rawData {
		if i == 0 {
			minVal = val
			maxVal = val
		} else {
			if val < minVal {
				minVal = val
			}
			if val > maxVal {
				maxVal = val
			}
		}
	}

	// figure out bin midpoints
	step := (maxVal - minVal) / float64(numBins)
	for i := 0; i < numBins; i++ {
		midpoints[i] = step * (float64(i) + 0.5)
	}
	// and now add to the bins
	for _, val := range rawData {
		whichBin := int(math.Floor((val - minVal) / step))
		if whichBin < 0 {
			whichBin = 0
		} else if whichBin >= numBins {
			whichBin = numBins - 1
		}
		vals[whichBin]++
	}
	return midpoints, vals, nil
}

func graphVals(vals []float64) {
	// for help in normalizing
	maxVal := float64(0)
	for _, val := range vals {
		if val > maxVal {
			maxVal = val
		}
	}
	if maxVal == 0 {
		fmt.Printf("no counts!\n")
		return
	}

	// precompute our scaled values just to save some repeated
	// arithmetic
	scaledVals := make([]int, len(vals))
	for i, val := range vals {
		scaledVals[i] = int(math.Floor((val / maxVal) * height))
	}

	colsPerBin := width / len(vals)
	for y := height; y > 0; y-- {
		for x := 0; x < width; x++ {
			whichBin := x / colsPerBin
			if y < scaledVals[whichBin] {
				fmt.Printf("#")
			} else {
				fmt.Printf(" ")
			}
		}
		fmt.Printf("\n")
	}
}

func readRawData() ([]float64, error) {
	result := make([]float64, 0)
	stdinReader := bufio.NewScanner(os.Stdin)
	for stdinReader.Scan() {
		thisLine := stdinReader.Text()
		val, err := strconv.ParseFloat(thisLine, 64)
		if err != nil {
			return result, fmt.Errorf("Error parsing line %s: %v", thisLine, err)
		}
		result = append(result, val)
	}
	return result, nil
}
