package arr

import "strconv"

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
