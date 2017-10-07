package utils

import "strings"

func Int32Ptr(i int32) *int32 { return &i }

func DotStr2DashStr(dotStr string) string {
	return strings.Replace(dotStr, ".", "-", -1)
}
