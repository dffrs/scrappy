package scraper

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

	if strings.Contains(match, ".") && strings.Contains(match, ",") {
		if strings.LastIndex(match, ".") > strings.LastIndex(match, ",") {
			match = strings.ReplaceAll(match, ",", "")
		} else {
			match = strings.ReplaceAll(match, ".", "")
			match = strings.ReplaceAll(match, ",", ".")
		}
	} else if strings.Contains(match, ",") {
		parts := strings.Split(match, ",")
		if len(parts[len(parts)-1]) == 2 {
			match = strings.ReplaceAll(match, ",", ".")
		} else {
			match = strings.ReplaceAll(match, ",", "")
		}
	}

	val, err := strconv.ParseFloat(match, 32)
	if err != nil {
		return 0, err
	}

	return float32(val), nil
}
