package http

import (
	"fmt"
	"testing"
	"time"

	"github.com/snowlyg/iris-admin/server/web"
	"github.com/snowlyg/iris-admin/server/web/web_gin"
)

var (
	timeover int64 = 5
	timeout  int64 = 10
)

func Test_Get(t *testing.T) {
	defer web_gin.Remove()
	web_gin.CONFIG.System.CacheType = "local"
	web_gin.CONFIG.System.Addr = "127.0.0.1:18088"
	go func() {
		web.Start(web_gin.Init())
	}()

	time.Sleep(3 * time.Second)

	client := NewClient(&Config{
		TimeOver: timeover,
		TimeOut:  timeout,
	})
	data := map[string]interface{}{}
	serviceResponseRestful := &ServerResponse{
		Path:     fmt.Sprintf("http://%s", web_gin.CONFIG.System.Addr),
		BaseAuth: true,
		Data:     data,
	}

	tests := []struct {
		name         string
		responseInfo *ServerResponse
	}{
		{
			name:         "获取接口列表",
			responseInfo: serviceResponseRestful,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := client.Get(tt.responseInfo)
			if err != nil {
				t.Errorf("Get get error = %v", err)
				return
			}
			if b == nil {
				t.Errorf("Get get got nil")
			}
		})
	}
}

func Test_Post(t *testing.T) {
	defer web_gin.Remove()
	web_gin.CONFIG.System.CacheType = "local"
	web_gin.CONFIG.System.Addr = "127.0.0.1:18088"
	go func() {
		web.Start(web_gin.Init())
	}()

	time.Sleep(3 * time.Second)

	client := NewClient(&Config{
		TimeOver: timeover,
		TimeOut:  timeout,
	})

	data := map[string]interface{}{}
	serviceResponseService := &ServerResponse{
		Path:     fmt.Sprintf("http://%s", web_gin.CONFIG.System.Addr),
		BaseAuth: false,
		Data:     data,
	}

	tests := []struct {
		name         string
		responseInfo *ServerResponse
		data         string
	}{
		{
			name:         "同步服务列表故障数据",
			responseInfo: serviceResponseService,
			data:         "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.Post(tt.responseInfo, tt.data)
			if err != nil {
				t.Errorf("Post get error = %v", err)
				return
			}
		})
	}
}

func Test_NewClient(t *testing.T) {
	args := []*Config{
		{
			TimeOver: 5,
			TimeOut:  5,
			Pwd:      "",  // driver redis password
			Headers:  nil, // request headers
		},
	}
	for _, arg := range args {
		t.Run("new client", func(t *testing.T) {
			client := NewClient(arg)
			if client == nil {
				t.Error("Get client is nil")
			}
		})
	}

}
