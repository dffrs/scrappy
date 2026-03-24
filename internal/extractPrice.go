package internal

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func extractPrice(input string) (float32, error) {
	re := regexp.MustCompile(`[\d.,]+`)
	match := re.FindString(input)

	if match == "" {
		return 0, fmt.Errorf("no number found")
	}

	// Case 1: both . and , → decide format
	if strings.Contains(match, ".") && strings.Contains(match, ",") {
		if strings.LastIndex(match, ".") > strings.LastIndex(match, ",") {
			// US format: 1,299.99
			match = strings.ReplaceAll(match, ",", "")
		} else {
			// EU format: 1.299,99
			match = strings.ReplaceAll(match, ".", "")
			match = strings.ReplaceAll(match, ",", ".")
		}
	} else if strings.Contains(match, ",") {
		// Could be decimal or thousands
		parts := strings.Split(match, ",")
		if len(parts[len(parts)-1]) == 2 {
			// decimal: 999,50
			match = strings.ReplaceAll(match, ",", ".")
		} else {
			// thousands: 1,299
			match = strings.ReplaceAll(match, ",", "")
		}
	}

	val, err := strconv.ParseFloat(match, 32)
	if err != nil {
		return 0, err
	}

	return float32(val), nil
}
