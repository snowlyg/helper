package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/snowlyg/helper/str"
)

var ErrInvalidDataStruct = errors.New("invalid response data struct")
var ErrBaseAuthConfig = errors.New("base auth config is error")
var ErrEmptyFileNameField = errors.New("must set field filename")

type client struct {
	config *Config
	cookie *http.Cookie
}

type Config struct {
	TimeOver   int64
	TimeOut    int64
	Headers    map[string]string // request headers
	CookieName string
	Host       string
}

// BaseAuth
type BaseAuth struct {
	Enable  bool
	Account string
	Pwd     string
}

// NewClient
func NewClient(configs ...*Config) *client {
	var config *Config
	if len(configs) == 0 || configs[0] == nil {
		config = &Config{TimeOut: 30, TimeOver: 5, Headers: map[string]string{}}
	} else {
		config = configs[0]
	}
	if config.TimeOut == 0 {
		config.TimeOut = 30
	}
	if config.TimeOver == 0 {
		config.TimeOver = 5
	}
	if config.Host == "" {
		config.Host = "http://127.0.0.1:7777"
	}
	return &client{config: config}
}

// NewResponse
func NewResponse(path string) *ServerResponse {
	return &ServerResponse{path: path}
}

type ServerResponse struct {
	path     string
	baseAuth *BaseAuth
	Data     interface{} `json:"data"`
	body     io.Reader
	fields   map[string]string
}

// SetBaseAuth
func (sr *ServerResponse) SetBaseAuth(account, pwd string) {
	sr.baseAuth = &BaseAuth{
		Enable:  true,
		Account: account,
		Pwd:     pwd,
	}
}

// BaseAuth
func (sr *ServerResponse) BaseAuth() *BaseAuth {
	return sr.baseAuth
}

// SetFields
func (sr *ServerResponse) SetFields(fields map[string]string) {
	sr.fields = fields
}

// GetFields
func (sr *ServerResponse) GetFields() map[string]string {
	return sr.fields
}

// SetUploadFile
func (sr *ServerResponse) SetUploadFile(name string) {
	f, err := os.Open(name)
	if err != nil {
		return
	}
	sr.body = f
}

// Close
func (sr *ServerResponse) Close() {
	if f, ok := sr.body.(*os.File); ok && f != nil {
		f.Close()
	}
}

func (n *client) GetCookie() *http.Cookie {
	return n.cookie
}

// Check
func (n *client) Check(sr *ServerResponse) error {
	ba := sr.BaseAuth()
	if ba == nil {
		return nil
	}
	if ba.Enable && (ba.Account == "" || ba.Pwd == "") {
		return ErrBaseAuthConfig
	}
	return nil
}

// getFullPath
func (n *client) getFullPath(path string) string {
	log.Println("fullpath:", path)
	return str.Join(n.config.Host, path)
}

// Post
func (n *client) Post(sr *ServerResponse, data string) ([]byte, error) {
	err := n.Check(sr)
	if err != nil {
		return nil, err
	}
	result := n.request("POST", n.getFullPath(sr.path), sr.BaseAuth(), strings.NewReader(data))
	if len(result) == 0 {
		return result, fmt.Errorf("Post %s 没有返回数据", n.getFullPath(sr.path))
	}
	if !json.Valid(result) || sr.Data == nil {
		sr.Data = string(result)
		return result, nil
	}
	err = json.Unmarshal(result, sr.Data)
	if err != nil {
		return result, fmt.Errorf("执行解码失败: %s 错误：%w ,结果: %v", n.getFullPath(sr.path), err, string(result))
	}
	return result, nil
}

// GetFile
func (n *client) GetFile(sr *ServerResponse) error {
	err := n.Check(sr)
	if err != nil {
		return err
	}
	_ = n.request("GET", n.getFullPath(sr.path), sr.BaseAuth(), nil)
	return nil
}

// Upload  上传文件
func (n *client) Upload(sr *ServerResponse) ([]byte, error) {
	err := n.Check(sr)
	if err != nil {
		return nil, err
	}

	if sr.fields == nil {
		return nil, ErrEmptyFileNameField
	}
	if filename, ok := sr.fields["filename"]; !ok || filename == "" {
		return nil, ErrEmptyFileNameField
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("file", sr.fields["filename"])
	if err != nil {
		return nil, fmt.Errorf("create form file %v", err)
	}

	_, err = io.Copy(fw, sr.body)
	if err != nil {
		return nil, fmt.Errorf("copying fileWriter %v", err)
	}

	for k, v := range sr.fields {
		_ = writer.WriteField(k, v)
	}

	err = writer.Close() // close writer before POST request
	if err != nil {
		return nil, fmt.Errorf("writerClose: %v", err)
	}

	n.config.Headers["Content-Type"] = writer.FormDataContentType()

	result := n.request("POST", n.getFullPath(sr.path), sr.BaseAuth(), body)
	if len(result) == 0 {
		return result, fmt.Errorf("Get %s 没有返回数据", n.getFullPath(sr.path))
	}

	if !json.Valid(result) || sr.Data == nil {
		sr.Data = string(result)
		return result, nil
	}

	err = json.Unmarshal(result, sr.Data)
	if err != nil {
		return result, fmt.Errorf("执行解码失败: %s 错误：%w ,结果: %v", n.getFullPath(sr.path), err, string(result))
	}

	return result, nil
}

// Get  获取数据
func (n *client) Get(sr *ServerResponse) ([]byte, error) {
	err := n.Check(sr)
	if err != nil {
		return nil, err
	}
	result := n.request("GET", n.getFullPath(sr.path), sr.BaseAuth(), nil)
	if len(result) == 0 {
		return result, fmt.Errorf("Get %s 没有返回数据", n.getFullPath(sr.path))
	}
	if !json.Valid(result) || sr.Data == nil {
		sr.Data = string(result)
		return result, nil
	}
	err = json.Unmarshal(result, sr.Data)
	if err != nil {
		return result, fmt.Errorf("执行解码失败: %s 获取服务解析返回内容报错 %w : ,结果:[%s]", n.getFullPath(sr.path), err, string(result))
	}
	return result, nil
}

func (n *client) request(method, fullpath string, ba *BaseAuth, body io.Reader) []byte {
	result := make(chan []byte, 30)
	T := time.NewTicker(time.Duration(n.config.TimeOver) * time.Second)
	go func() {
		t := time.Duration(n.config.TimeOut) * time.Second
		Client := http.Client{Timeout: t}
		req, err := http.NewRequest(method, fullpath, body)
		if err != nil {
			result <- nil
			return
		}

		if len(n.config.Headers) > 0 {
			for key, value := range n.config.Headers {
				req.Header.Set(key, value)
			}
		} else {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
		}
		if ba != nil && ba.Enable {
			req.SetBasicAuth(ba.Account, ba.Pwd)
		}
		var resp *http.Response
		resp, err = Client.Do(req)
		if err != nil {
			result <- nil
			return
		}
		if n.config.CookieName != "" && resp.Cookies() != nil {
			for _, cookie := range resp.Cookies() {
				if cookie.Name == n.config.CookieName {
					n.cookie = cookie
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
