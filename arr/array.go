package arr

import (
	"strconv"
)

//InArray
func InArray[T comparable](items []T, s T) bool {
	for _, item := range items {
		if item == s {
			return true
		}
	}
	return false
}

// 连接 unit slice 为字符串
func UnitJoin(ss []uint, sep string) string {
	var rs string
	for index, item := range ss {
		itemS := strconv.FormatUint(uint64(item), 10)
		if index < len(ss)-1 {
			rs += itemS + sep
		} else {
			rs += itemS
		}
	}
	return rs
}
