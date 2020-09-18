package tools

import "strings"

// 替换时间的T和Z
func ReplaceTime(t string) string {
	t = strings.Replace(t, "T", " ", -1)
	t = strings.Replace(t, "Z", "", -1)
	t = strings.Replace(t, "+0800 CS", "", -1)
	ts := strings.Split(t, ".")
	return ts[0]
}
