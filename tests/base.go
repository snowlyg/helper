package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/snowlyg/helper/str"
)

type Client struct {
	baseURL string
	t       *testing.T
	handler http.Handler
}

func New(url string, t *testing.T, handler http.Handler) *Client {
	return &Client{
		baseURL: url,
		t:       t,
		handler: handler,
	}
}

func (c *Client) Expect() *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL: c.baseURL,
		Client: &http.Client{
			Transport: httpexpect.NewBinder(c.handler),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(c.t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(c.t, true),
			httpexpect.NewCurlPrinter(c.t),
			httpexpect.NewCompactPrinter(c.t),
		},
	})
}

func (c *Client) Login(url string, res Responses, datas ...map[string]interface{}) *httpexpect.Expect {
	data := LoginParams
	if len(datas) > 0 {
		data = datas[0]
	}
	obj := c.Expect().POST(url).WithJSON(data).Expect().Status(http.StatusOK).JSON().Object()
	if res == nil {
		res = LoginResponse
	}
	token := res.Test(obj).GetString("AccessToken")
	fmt.Printf("access_token is '%s'\n", token)
	return c.Expect().Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", str.Join("Bearer ", token))
	})
}

func (c *Client) Logout(url string, res Responses) {
	if res == nil {
		res = LogoutResponse
	}
	obj := c.Expect().GET(url).Expect().Status(http.StatusOK).JSON().Object()
	res.Test(obj)
}
