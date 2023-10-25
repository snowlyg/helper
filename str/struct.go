package str

import (
	"reflect"
	"strings"
)

// StructToMap 利用反射将结构体转化为map
func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}

// 连接字符串
func Join(strs ...string) string {
	var builder strings.Builder
	if len(strs) == 0 {
		return ""
	}
	for _, str := range strs {
		builder.WriteString(str)
	}
	return builder.String()
}
