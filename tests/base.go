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

func (c *Client) Login(url string, res Responses, datas ...map[string]interface{}) error {
	data := LoginParams
	if len(datas) > 0 {
		data = datas[0]
	}
	if res == nil {
		res = LoginResponse
	}
	token := c.POST(url, res, data).GetString("data.accessToken")
	fmt.Printf("access_token is '%s'\n", token)
	if token == "" {
		return fmt.Errorf("access_token is empty")
	}
	c.expect = c.expect.Builder(func(req *httpexpect.Request) {
		req.WithHeader("Authorization", str.Join("Bearer ", token))
	})
	return nil
}

func (c *Client) Logout(url string, res Responses) {
	if res == nil {
		res = LogoutResponse
	}
	c.GET(url, res)
}

// POST
func (c *Client) POST(url string, res Responses, data map[string]interface{}) Responses {
	obj := c.expect.POST(url).WithJSON(data).Expect().Status(http.StatusOK).JSON().Object()
	return res.Test(obj)
}

// GET
func (c *Client) GET(url string, res Responses, datas ...map[string]interface{}) Responses {
	req := c.expect.GET(url)
	if len(datas) > 0 {
		req = req.WithQueryObject(datas[0])
	}
	obj := req.Expect().Status(http.StatusOK).JSON().Object()
	return res.Test(obj)
}
// DELETE
func (c *Client) DELETE(url string, res Responses, datas ...map[string]interface{}) Responses {
	req := c.expect.DELETE(url)
	if len(datas) > 0 {
		req = req.WithQueryObject(datas[0])
	}
	obj := req.Expect().Status(http.StatusOK).JSON().Object()
	return res.Test(obj)
}
