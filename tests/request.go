package tests

var (
	RequestParams = map[string]interface{}{"page": 1, "pageSize": 10}
	LoginParams   = map[string]interface{}{"username": "admin", "password": "123456"}
	LoginResponse = Responses{
		{Key: "code", Value: 2000},
		{Key: "message", Value: "请求成功"},
		{Key: "data",
			Value: Responses{
				{Key: "accessToken", Value: "", Type: "notempty"},
			},
		},
	}
	LogoutResponse = Responses{
		{Key: "code", Value: 2000},
		{Key: "message", Value: "请求成功"},
	}
	ResponseDatas = Responses{
		{Key: "pageSize", Value: 10},
		{Key: "page", Value: 1},
	}
)

// func PostList(auth *httpexpect.Expect, url string, reqs Responses, status int, message string, res ...map[string]interface{}) {
// 	obj := auth.POST(url).WithJSON(res).Expect().Status(http.StatusOK).JSON().Object()
// 	if reqs == nil {
// 		return
// 	}
// 	reqs.Test(obj)
// }

// func GetList(auth *httpexpect.Expect, url string, status int, message string, reqs Responses) {
// 	obj := auth.GET(url).Expect().Status(http.StatusOK).JSON().Object()
// 	if reqs == nil {
// 		return
// 	}
// 	reqs.Test(obj)
// }

// func Create(auth *httpexpect.Expect, url string, create map[string]interface{}, keys Responses, status int, message string) {
// 	obj := auth.POST(url).WithJSON(create).Expect().Status(http.StatusOK).JSON().Object()
// 	if status != http.StatusOK {
// 		return
// 	}
// 	if keys == nil {
// 		return
// 	}
// 	keys.Scan(obj)
// }

// func Update(auth *httpexpect.Expect, url string, update map[string]interface{}, status int, message string) {
// 	obj := auth.PUT(url).WithJSON(update).Expect().Status(http.StatusOK).JSON().Object()
// 	obj.Keys().ContainsOnly("status", "data", "message")
// 	obj.Value("status").Number().Equal(status)
// 	obj.Value("message").String().Equal(message)
// }

// func Get(auth *httpexpect.Expect, url string, query map[string]interface{}, status int, message string, keys ...Responses) {
// 	req := auth.GET(url)
// 	if query != nil {
// 		req = req.WithQueryObject(query)
// 	}
// 	obj := req.Expect().Status(status).JSON().Object()
// 	obj.Keys().ContainsOnly("status", "data", "message")
// 	obj.Value("status").Number().Equal(status)
// 	obj.Value("message").String().Equal(message)
// 	if len(keys) == 0 {
// 		return
// 	}
// 	//返回单个数据
// 	if len(keys) == 1 {
// 		keys[0].Test(obj.Value("data").Object())
// 		return
// 	}
// 	// 返回数组数据
// 	for m, ks := range keys {
// 		if ks == nil {
// 			return
// 		}
// 		ks.Test(obj.Value("data").Array().Element(m).Object())
// 	}
// }

// func ScanById(auth *httpexpect.Expect, url string, id uint, query map[string]interface{}, keys Responses, status int, message string) {
// 	req := auth.GET(url)
// 	if query != nil {
// 		req = req.WithQueryObject(query)
// 	}
// 	obj := req.Expect().Status(http.StatusOK).JSON().Object()
// 	obj.Keys().ContainsOnly("status", "data", "message")
// 	obj.Value("status").Number().Equal(status)
// 	obj.Value("message").String().Equal(message)
// 	if status == http.StatusOK {
// 		if keys != nil {
// 			keys.Scan(obj.Value("data").Object())
// 		}
// 	}
// }

// func Post(auth *httpexpect.Expect, url string, data map[string]interface{}, status int, message string) {
// 	obj := auth.POST(url).WithJSON(data).Expect().Status(http.StatusOK).JSON().Object()
// 	obj.Keys().ContainsOnly("status", "data", "message")
// 	obj.Value("status").Number().Equal(status)
// 	obj.Value("message").String().Equal(message)
// }

// func Delete(auth *httpexpect.Expect, url string, status int, message string) {
// 	obj := auth.DELETE(url).Expect().Status(http.StatusOK).JSON().Object()
// 	obj.Keys().ContainsOnly("status", "data", "message")
// 	obj.Value("status").Number().Equal(status)
// 	obj.Value("message").String().Equal(message)
// }

// func DeleteMutil(auth *httpexpect.Expect, url string, data map[string]interface{}, status int, message string) {
// 	obj := auth.DELETE(url).WithJSON(data).Expect().Status(http.StatusOK).JSON().Object()
// 	obj.Keys().ContainsOnly("status", "data", "message")
// 	obj.Value("status").Number().Equal(status)
// 	obj.Value("message").String().Equal(message)
// }
