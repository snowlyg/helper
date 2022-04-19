package arr

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

//InArray
func InArray[T constraints.Ordered](items []T, s T) bool {
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
