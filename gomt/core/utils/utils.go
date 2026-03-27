package utils

import (
	"strings"

	"github.com/spf13/cast"
)

func CompareVersion(version, version2 string) int {
	args := strings.Split(version, ".")
	args2 := strings.Split(version, ".")
	size := len(args)
	size2 := len(args2)
	max := size
	if size2 > size {
		max = size2
	}

	for i := 0; i < max; i++ {
		v1 := 0
		v2 := 0
		if i < size {
			v1 = cast.ToInt(args[i])
		}
		if i < size2 {
			v2 = cast.ToInt(args2[i])
		}
		if v1 == v2 {
			continue
		} else if v1 > v2 {
			return 1
		} else {
			return -1
		}
	}
	return 0
}
