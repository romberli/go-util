package linux

import (
	"bufio"
	"os"

	"github.com/pingcap/errors"
)

const (
	DefaultEstimateLineSize = 1024
	MinStartPosition        = 0
)

// TailN try get the latest n line of the file.
func TailN(fileName string, n int) (lines []string, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.AddStack(err)
	}
	defer func() { _ = file.Close() }()

	estimateLineSize := DefaultEstimateLineSize

	stat, err := os.Stat(fileName)
	if err != nil {
		return nil, errors.AddStack(err)
	}

	start := int(stat.Size()) - n*estimateLineSize
	if start < MinStartPosition {
		start = MinStartPosition
	}

	_, err = file.Seek(int64(start), MinStartPosition /*means relative to the origin of the file*/)
	if err != nil {
		return nil, errors.AddStack(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}

	return
}
