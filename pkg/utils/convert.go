package utils

import "strconv"

// The vendor SQL Account API represents decimal/numeric fields
// inconsistently as either JSON numbers or numeric strings, so these
// conversions accept either when parsing its responses.

func AnyToInt(v any) int {
	switch t := v.(type) {
	case float64:
		return int(t)
	case string:
		n, _ := strconv.Atoi(t)
		return n
	default:
		return 0
	}
}

func AnyToString(v any) string {
	s, _ := v.(string)
	return s
}

func AnyToStringPtr(v any) *string {
	s, ok := v.(string)
	if !ok || s == "" {
		return nil
	}
	return &s
}

func AnyToFloatPtr(v any) *float64 {
	switch t := v.(type) {
	case float64:
		return &t
	case string:
		f, err := strconv.ParseFloat(t, 64)
		if err != nil {
			return nil
		}
		return &f
	default:
		return nil
	}
}
