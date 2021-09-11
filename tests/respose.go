package tests

import (
	"reflect"
	"sync"

	"github.com/gavv/httpexpect/v2"
)

var wg sync.WaitGroup

// type Param struct {
// 	Name string
// 	Args map[string]interface{}
// 	Responses
// }

type Responses []Response
type Response struct {
	Type  string
	Key   string
	Value interface{}
}

// Keys
func (res Responses) Keys() []string {
	keys := []string{}
	for _, re := range res {
		keys = append(keys, re.Key)
	}
	return keys
}

func IdKeys() Responses {
	return Responses{
		{Key: "id", Value: uint(0)},
	}
}

func (res Responses) Test(object *httpexpect.Object) Responses {
	for _, rs := range res {
		wg.Add(1)
		if rs.Value == nil {
			wg.Done()
			continue
		}
		switch reflect.TypeOf(rs.Value).String() {
		case "string":
			if rs.Value.(string) == "notempty" {
				object.Value(rs.Key).String().NotEmpty()
			} else {
				object.Value(rs.Key).String().Equal(rs.Value.(string))
			}
		case "float64":
			object.Value(rs.Key).Number().Equal(rs.Value.(float64))
		case "uint":
			object.Value(rs.Key).Number().Equal(rs.Value.(uint))
		case "int":
			object.Value(rs.Key).Number().Equal(rs.Value.(int))
		case "[]base.ResponseKeys":
			object.Value(rs.Key).Array().Length().Equal(len(rs.Value.([]Responses)))
			length := int(object.Value(rs.Key).Array().Length().Raw())
			if length > 0 && len(rs.Value.([]Responses)) == length {
				for i := 0; i < length; i++ {
					rs.Value.([]Responses)[i].Test(object.Value(rs.Key).Array().Element(i).Object())
				}
			}
		case "[]uint":
			object.Value(rs.Key).Array().Length().Equal(len(rs.Value.([]uint)))
			length := int(object.Value(rs.Key).Array().Length().Raw())
			if length > 0 && len(rs.Value.([]uint)) == length {
				for i := 0; i < length; i++ {
					object.Value(rs.Key).Array().Element(i).Number().Equal(rs.Value.([]uint)[i])
				}
			}
		case "[]string":
			object.Value(rs.Key).Array().Length().Equal(len(rs.Value.([]string)))
			length := int(object.Value(rs.Key).Array().Length().Raw())
			if length > 0 && len(rs.Value.([]string)) == length {
				for i := 0; i < length; i++ {
					object.Value(rs.Key).Array().Element(i).String().Equal(rs.Value.([]string)[i])
				}
			}

		default:
			wg.Done()
			continue
		}
		wg.Done()
	}
	wg.Wait()
	return res.Scan(object)
}

func (res Responses) Scan(object *httpexpect.Object) Responses {
	for k, rk := range res {
		wg.Add(1)
		if !Exist(object, rk.Key) {
			wg.Done()
			continue
		}
		switch reflect.TypeOf(rk.Value).String() {
		case "string":
			res[k].Value = object.Value(rk.Key).String().Raw()
		case "uint":
			res[k].Value = uint(object.Value(rk.Key).Number().Raw())
		case "int":
			res[k].Value = int(object.Value(rk.Key).Number().Raw())
		case "int32":
			res[k].Value = int32(object.Value(rk.Key).Number().Raw())
		case "float64":
			res[k].Value = object.Value(rk.Key).Number().Raw()
		case "[]base.ResponseKeys":
			object.Value(rk.Key).Array().Length().Equal(len(rk.Value.([]Responses)))
			length := int(object.Value(rk.Key).Array().Length().Raw())
			if length > 0 && len(rk.Value.([]Responses)) == length {
				for i := 0; i < length; i++ {
					rk.Value.([]Responses)[i].Scan(object.Value(rk.Key).Array().Element(i).Object())
				}
			}
		case "[]string":
			length := int(object.Value(rk.Key).Array().Length().Raw())

			if length == 0 {
				wg.Done()
				continue
			}
			reskey, ok := res[k].Value.([]string)
			if ok {
				var strings []string
				for i := 0; i < length; i++ {
					strings = append(reskey, object.Value(rk.Key).Array().Element(i).String().Raw())
				}
				res[k].Value = strings
			}
		default:
			wg.Done()
			continue
		}
		wg.Done()
	}
	wg.Wait()
	return res
}

func Exist(object *httpexpect.Object, key string) bool {
	objectKyes := object.Keys().Raw()
	for _, objectKey := range objectKyes {
		if key == objectKey.(string) {
			return true
		}
	}
	return false
}

func (rks Responses) GetString(key string) string {
	for _, rk := range rks {
		if key == rk.Key {
			if rk.Value == nil {
				return ""
			}
			switch reflect.TypeOf(rk.Value).String() {
			case "string":
				return rk.Value.(string)
			}
		}
	}
	return ""
}

func (rks Responses) GetStrArray(key string) []string {
	for _, rk := range rks {
		if key == rk.Key {
			if rk.Value == nil {
				return nil
			}
			switch reflect.TypeOf(rk.Value).String() {
			case "[]string":
				return rk.Value.([]string)
			}
		}
	}
	return nil
}

func (rks Responses) GetResponses(key string) []Responses {
	for _, rk := range rks {
		if key == rk.Key {
			if rk.Value == nil {
				return nil
			}
			switch reflect.TypeOf(rk.Value).String() {
			case "[]base.ResponseKeys":
				return rk.Value.([]Responses)
			}
		}
	}
	return nil
}

func (rks Responses) GetUint(key string) uint {
	for _, rk := range rks {
		if key == rk.Key {
			if rk.Value == nil {
				return 0
			}
			switch reflect.TypeOf(rk.Value).String() {
			case "float64":
				return uint(rk.Value.(float64))
			case "int32":
				return uint(rk.Value.(int32))
			case "uint":
				return rk.Value.(uint)
			case "int":
				return uint(rk.Value.(int))
			}
		}
	}
	return 0
}

func (rks Responses) GetInt(key string) int {
	for _, rk := range rks {
		if key == rk.Key {
			if rk.Value == nil {
				return 0
			}
			switch reflect.TypeOf(rk.Value).String() {
			case "float64":
				return int(rk.Value.(float64))
			case "int":
				return rk.Value.(int)
			case "int32":
				return int(rk.Value.(int32))
			case "uint":
				return int(rk.Value.(uint))
			}
		}
	}
	return 0
}
func (rks Responses) GetInt32(key string) int32 {
	for _, rk := range rks {
		if key == rk.Key {
			if rk.Value == nil {
				return 0
			}
			switch reflect.TypeOf(rk.Value).String() {
			case "float64":
				return int32(rk.Value.(float64))
			case "int32":
				return rk.Value.(int32)
			case "int":
				return int32(rk.Value.(int))
			case "uint":
				return int32(rk.Value.(uint))
			}
		}
	}
	return 0
}

func (res Responses) GetId() uint {
	return res.GetUint("id")
}
