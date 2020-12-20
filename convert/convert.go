package convert

import "strconv"

// string转换成int64
func StringMustToInt64(v string) int64 {
	if v == "" {
		return 0
	}
	i, _ := strconv.ParseInt(v, 10, 64)

	return i
}
