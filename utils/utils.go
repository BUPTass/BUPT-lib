package utils

import "strconv"

func CastInt(in string) int {
	i, _ := strconv.Atoi(in)
	return i
}
