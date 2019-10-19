package scanner

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

// DependencyScanner reads co-change MDG files
type DependencyScanner struct {
	scanner *bufio.Scanner
}

// NewDependencyScanner returns a new DependencyScanner to read from r
func NewDependencyScanner(r io.Reader) *DependencyScanner {
	return &DependencyScanner{bufio.NewScanner(r)}
}

// Scan advances the scanner to the next dependency
func (ds *DependencyScanner) Scan() bool {
	for ds.scanner.Scan() {
		if strings.TrimSpace(ds.scanner.Text()) != "" {
			return true
		}
	}
	return false
}

// Err return the first non-EOF error that was encountered by the scanner
func (ds *DependencyScanner) Err() error {
	return ds.scanner.Err()
}

// Dependency returns the most recent dependency generated by a call to Scan
func (ds *DependencyScanner) Dependency() struct {
	From         []string
	To           string
	SupportCount int
	Confidence   float64
	CommitsCount int
} {
	arr := strings.Split(strings.TrimSpace(ds.scanner.Text()), "\t")
	if len(arr) < 2 {
		arr = strings.Split(ds.scanner.Text(), " ")
	}
	var i int
	for i = len(arr) - 1; i > -1; i-- {
		_, err := strconv.ParseFloat(arr[i], 32)
		if err != nil {
			break
		}
	}
	entities := arr[0 : i+1]
	var numbers []string
	if i < len(arr)-1 {
		numbers = arr[i+1:]
	}
	supportCount := 0
	if len(numbers) > 0 {
		supportCount, _ = strconv.Atoi(numbers[0])
	}
	confidence := 0.0
	if len(numbers) > 1 {
		confidence, _ = strconv.ParseFloat(numbers[1], 32)
	}
	commitsCount := 0
	if len(numbers) > 3 {
		commitsCount, _ = strconv.Atoi(numbers[3])
	}
	return struct {
		From         []string
		To           string
		SupportCount int
		Confidence   float64
		CommitsCount int
	}{
		entities[:len(entities)-1], entities[len(entities)-1],
		supportCount,
		confidence,
		commitsCount,
	}
}
