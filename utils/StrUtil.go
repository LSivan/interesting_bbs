package utils

import (
	"strconv"
	"strings"
	"encoding/json"
)

func NoHtml(str string) string {
	return strings.Replace(strings.Replace(str, "<script", "&lt;script", -1), "script>", "script&gt;", -1)
}

func MustInt(str string) int {
	if v, err := strconv.Atoi(str); err == nil {
		return v
	}
	return 0
}

func MustJson(data interface{}) string {
	s, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(s)
}