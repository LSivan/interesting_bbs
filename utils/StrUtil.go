package utils

import (
	"strings"
	"strconv"
)

func NoHtml(str string) string {
	return strings.Replace(strings.Replace(str, "<script", "&lt;script", -1), "script>", "script&gt;", -1)
}

func MustInt(str string) int {
	if v,err := strconv.Atoi(str);err == nil {
		return v
	}
	return 0
}