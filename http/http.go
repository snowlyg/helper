package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

type client struct {
	Config *Config
	Cookie *http.Cookie
}

type Config struct {
	TimeOver   int64
	TimeOut    int64
	Account    string            // 账号
	Pwd        string            // 密码
	Headers    map[string]string // request headers
	CookieName string
}

func NewClient(config *Config) *client {
	if config.TimeOut == 0 {
		config.TimeOut = 30
	}
	if config.TimeOver == 0 {
		config.TimeOver = 5
	}
	return &client{Config: config}
}

type ServerResponse struct {
	Path     string      `json:"path"`
	BaseAuth bool        `json:"baseAuth"`
	Data     interface{} `json:"data"`
	Body     io.Reader
	Fields   map[string]string
}

// Post  提交数据
func (n *client) Post(sr *ServerResponse, data string) ([]byte, error) {
	result := n.request("POST", sr.Path, sr.BaseAuth, strings.NewReader(data))
	if len(result) == 0 {
		return result, fmt.Errorf("Post %s 没有返回数据", sr.Path)
	}
	if !json.Valid(result) {
		return result, fmt.Errorf("result is not valid")
	}
	err := json.Unmarshal(result, sr.Data)
	if err != nil {
		return result, fmt.Errorf("执行解码失败: %s 错误：%w ,结果: %v", sr.Path, err, string(result))
	}
	return result, nil
}

// GetFile  下载文件
func (n *client) GetFile(sr *ServerResponse) ([]byte, error) {
	result := n.request("GET", sr.Path, sr.BaseAuth, nil)
	if len(result) == 0 {
		return result, fmt.Errorf("Get %s 没有返回数据", sr.Path)
	}
	if !json.Valid(result) {
		return result, fmt.Errorf("result is not valid")
	}
	err := json.Unmarshal(result, sr.Data)
	if err != nil {
		return result, fmt.Errorf("执行解码失败: %s 错误：%w ,结果: %v", sr.Path, err, string(result))
	}
	return result, nil
}

// Upload  上传文件
func (n *client) Upload(sr *ServerResponse) ([]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("file", sr.Fields["filename"])
	if err != nil {
		return nil, fmt.Errorf("CreateFormFile %v", err)
	}

	_, err = io.Copy(fw, sr.Body)
	if err != nil {
		return nil, fmt.Errorf("copying fileWriter %v", err)
	}

	for k, v := range sr.Fields {
		_ = writer.WriteField(k, v)
	}
	err = writer.Close() // close writer before POST request
	if err != nil {
		return nil, fmt.Errorf("writerClose: %v", err)
	}
	n.Config.Headers["Content-Type"] = writer.FormDataContentType()
	result := n.request("POST", sr.Path, sr.BaseAuth, body)
	if len(result) == 0 {
		return result, fmt.Errorf("Get %s 没有返回数据", sr.Path)
	}
	if !json.Valid(result) {
		return result, fmt.Errorf("result is not valid")
	}
	err = json.Unmarshal(result, sr.Data)
	if err != nil {
		return result, fmt.Errorf("执行解码失败: %s 错误：%w ,结果: %v", sr.Path, err, string(result))
	}
	return result, nil
}

// Get  获取数据
func (n *client) Get(sr *ServerResponse) ([]byte, error) {

	result := n.request("GET", sr.Path, sr.BaseAuth, nil)
	if len(result) == 0 {
		return result, fmt.Errorf("Get %s 没有返回数据", sr.Path)
	}
	if !json.Valid(result) {
		return result, fmt.Errorf("result is not valid")
	}
	err := json.Unmarshal(result, sr.Data)
	if err != nil {
		return result, fmt.Errorf("执行解码失败: %s 获取服务解析返回内容报错 %w : ,结果:[%s]", sr.Path, err, string(result))
	}

	return result, nil
}

func (n *client) request(method, url string, basicAuth bool, body io.Reader) []byte {
	result := make(chan []byte, 30)
	T := time.NewTicker(time.Duration(n.Config.TimeOver) * time.Second)
	go func() {
		t := time.Duration(n.Config.TimeOut) * time.Second
		Client := http.Client{Timeout: t}
		req, err := http.NewRequest(method, url, body)
		if err != nil {
			result <- nil
			return
		}

		if len(n.Config.Headers) > 0 {
			for key, value := range n.Config.Headers {
				req.Header.Set(key, value)
			}
		} else {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		}
		if basicAuth {
			req.SetBasicAuth(n.Config.Account, n.Config.Pwd)
		}
		var resp *http.Response
		resp, err = Client.Do(req)
		if err != nil {
			result <- nil
			return
		}
		if n.Config.CookieName != "" && resp.Cookies() != nil {
			for _, cookie := range resp.Cookies() {
				if cookie.Name == n.Config.CookieName {
					n.Cookie = cookie
				}
			}
		}
		defer resp.Body.Close()

		buf := bytes.NewBuffer(nil)
		io.Copy(buf, resp.Body)
		result <- buf.Bytes()

	}()

	for {
		select {
		case x := <-result:
			return x
		case <-T.C:
			return []byte("请求超时")
		}
	}
}
