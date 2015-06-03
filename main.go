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

// TODO: add log bins
var options struct {
	numBins    int
	minVal     float64
	maxVal     float64
	skipErrors bool
	skipNaN    bool
}

func init() {
	flag.IntVar(&options.numBins, "num-bins", defaultNumBins,
		"number of bins to include (legal values: 10, 20, 40, 80)")
	flag.Float64Var(&options.minVal, "min-val", math.Inf(1), "minimum value to graph")
	flag.Float64Var(&options.maxVal, "max-val", math.Inf(-1), "maximum value to graph")
	flag.BoolVar(&options.skipErrors, "skip-errors", false, "ignore lines with errors")
	flag.BoolVar(&options.skipNaN, "skip-nan", true,
		"if true, ignores lines that parse to NaN. otherwise, bails out with error")
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
		fmt.Fprintf(os.Stderr, "Illegal number of bins: %d\n", options.numBins)
		os.Exit(1)
	}

	rawData, err := readRawData()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading data: %v\n")
		os.Exit(1)
	}
	if len(rawData) == 0 {
		fmt.Fprintf(os.Stderr, "No data read!")
		return
	}

	rangeMin, rangeMax := getBinRange(rawData)
	if math.IsInf(rangeMin, 0) || math.IsInf(rangeMax, 0) {
		fmt.Fprintf(os.Stderr, "bad bin range: (%f, %f)", rangeMin, rangeMax)
		os.Exit(1)
	}

	vals := makeBins(rawData, rangeMin, rangeMax, options.numBins)
	drawHistogram(vals)
	drawXAxis(rangeMin, rangeMax)
}

// This is designed to **NOT** return NaNs under any circumstances
func readRawData() ([]float64, error) {
	result := make([]float64, 0)
	stdinReader := bufio.NewScanner(os.Stdin)
	for stdinReader.Scan() {
		thisLine := stdinReader.Text()
		val, err := strconv.ParseFloat(thisLine, 64)
		if err != nil {
			if options.skipErrors {
				continue
			}
			return result, fmt.Errorf("Error parsing line %s: %v", thisLine, err)
		}
		if math.IsNaN(val) {
			if options.skipNaN {
				continue
			}
			return result, fmt.Errorf("Found NaN parsing following line: %s", thisLine)
		}
		result = append(result, val)
	}
	return result, nil
}

func getBinRange(rawData []float64) (float64, float64) {
	// figure out range
	minVal := options.minVal
	maxVal := options.maxVal
	for _, val := range rawData {
		if val < minVal {
			minVal = val
		}
		if val > maxVal {
			maxVal = val
		}
	}

	return minVal, maxVal
}

func makeBins(rawData []float64, rangeMin, rangeMax float64, numBins int) []float64 {
	vals := make([]float64, numBins)
	step := (rangeMax - rangeMin) / float64(numBins)
	for _, val := range rawData {
		whichBin := int(math.Floor((val - rangeMin) / step))
		if whichBin < 0 {
			whichBin = 0
		} else if whichBin >= numBins {
			whichBin = numBins - 1
		}
		vals[whichBin]++
	}
	return vals
}

func drawHistogram(vals []float64) {
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
	total := 0
	scaledVals := make([]int, len(vals))
	for i, val := range vals {
		total += int(val)
		scaledVals[i] = int(math.Floor((val / maxVal) * height))
	}
	fmt.Printf("N = %d\n", total)

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

func drawXAxis(rangeMin, rangeMax float64) {
	tickLocations := []int{9, 29, 49, 69}
	column := 0
	for _, tickColumn := range tickLocations {
		for column < tickColumn {
			fmt.Printf("-")
			column++
		}
		fmt.Printf("+")
	}
	for column < width {
		fmt.Printf("-")
		column++
	}
	fmt.Printf("\n")

	// now plot the values
	column = 0
	step := (rangeMax - rangeMin) / float64(width)
	for _, tickColumn := range tickLocations {
		tickValue := rangeMin + (step * (float64(tickColumn) + 0.5))
		tickStr := strconv.FormatFloat(tickValue, 'g', 8, 64)
		startCol := tickColumn - (len(tickStr) / 2)
		for column < startCol {
			fmt.Printf(" ")
			column++
		}
		fmt.Printf(tickStr)
		column += len(tickStr)
		// TODO: figure out why this makes things look better :/
		column--
	}
	fmt.Printf("\n")
}
