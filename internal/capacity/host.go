package capacity

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

func hostCPUs() int {
	return runtime.NumCPU()
}

func parseMemTotalMiB(meminfo []byte) (int64, error) {
	lines := strings.Split(string(meminfo), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) < 2 {
				return 0, fmt.Errorf("host: invalid meminfo line: %s", line)
			}
			total, err := strconv.ParseInt(fields[1], 10, 64)
			if err != nil {
				return 0, fmt.Errorf("host: invalid meminfo line: %s", line)
			}
			return total / 1024, nil
		}
	}

	return 0, fmt.Errorf("host: MemTotal not found in meminfo")
}
