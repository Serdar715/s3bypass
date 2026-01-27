package utils

import (
	"strconv"
	"strings"
)

// ParseIntList parses a comma-separated string of integers into a map for fast lookup.
// Example: "200,403,500" -> map[200]struct{}{...}
func ParseIntList(input string) (map[int]struct{}, error) {
	result := make(map[int]struct{})
	if input == "" {
		return result, nil
	}

	parts := strings.Split(input, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		
		val, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		result[val] = struct{}{}
	}
	return result, nil
}
