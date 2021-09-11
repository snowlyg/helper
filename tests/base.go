package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/snowlyg/helper/str"
)

type Client struct {
	expect *httpexpect.Expect
}

func New(url string, t *testing.T, handler http.Handler) *Client {
	return &Client{
		expect: httpexpect.WithConfig(httpexpect.Config{
			BaseURL: url,
			Client: &http.Client{
				Transport: httpexpect.NewBinder(handler),
				Jar:       httpexpect.NewJar(),
			},
			Reporter: httpexpect.NewAssertReporter(t),
			Printers: []httpexpect.Printer{
				httpexpect.NewDebugPrinter(t, true),
				httpexpect.NewCurlPrinter(t),
				httpexpect.NewCompactPrinter(t),
			},
		}),
	}
}

func (c *Client) Login(url string, res Responses, datas ...map[string]interface{}) {
	data := LoginParams
	if len(datas) > 0 {
		data = datas[0]
	}
	obj := c.expect.POST(url).WithJSON(data).Expect().Status(http.StatusOK).JSON().Object()
	if res == nil {
		res = LoginResponse
	}
	token := res.Test(obj).GetString("data.accessToken")
	fmt.Printf("access_token is '%s'\n", token)
	c.expect = c.expect.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", str.Join("Bearer ", token))
	})
}

func (c *Client) Logout(url string, res Responses) {
	if res == nil {
		res = LogoutResponse
	}
	obj := c.expect.GET(url).Expect().Status(http.StatusOK).JSON().Object()
	res.Test(obj)
}
