package xstrings

import "strconv"

func Int64s2Strings(int64s []int64) []string {
	strings := make([]string, len(int64s))
	for i, v := range int64s {
		strings[i] = strconv.FormatInt(v, 10)
	}
	return strings
}

func Int32s2Strings(ints []int64) []string {
	strings := make([]string, len(ints))
	for i, v := range ints {
		strings[i] = strconv.FormatInt(v, 10)
	}
	return strings
}

func Ints2Strings(ints []int) []string {
	strings := make([]string, len(ints))
	for i, v := range ints {
		strings[i] = strconv.Itoa(v)
	}
	return strings
}
