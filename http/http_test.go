package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
)

func startGin() {
	r := gin.Default()
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	r.POST("/foo", func(c *gin.Context) {
		c.JSON(http.StatusOK, "foo")
	})

	r.GET("/bar", func(c *gin.Context) {
		c.JSON(http.StatusOK, "bar")
	})

	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})
	r.StaticFS("/txt", http.Dir("./txt"))

	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.POST("/upload", func(c *gin.Context) {
		// single file
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, "upload failed!")
			return
		}

		// Upload the file to specific dst.
		c.SaveUploadedFile(file, "./txt")

		c.JSON(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})

	// Listen and Server in http://0.0.0.0:7777
	r.Run(":7777")
}

func TestNewClient(t *testing.T) {

	client := NewClient()
	t.Run("test new client", func(t *testing.T) {
		if client == nil {
			t.Error("client is nil")
			return
		}
		if client.config.TimeOut != 30 {
			t.Errorf("client default timeout want %d but get %d", 30, client.config.TimeOut)
			return
		}
		if client.config.TimeOver != 5 {
			t.Errorf("client default timeover want %d but get %d", 5, client.config.TimeOver)
			return
		}
		if client.config.Host != "http://127.0.0.1:7777" {
			t.Errorf("client default timeover want %s but get %s", "http://127.0.0.1:7777", client.config.Host)
			return
		}
	})

	response := NewResponse("/foo")
	fullpath := client.getFullPath("/fullpath")
	if client.getFullPath("/fullpath") != "http://127.0.0.1:7777/fullpath" {
		t.Errorf("client default timeover want %s but get %s", "http://127.0.0.1:7777", fullpath)
		return
	}
	t.Run("test new response", func(t *testing.T) {
		if response == nil {
			t.Error("response is nil")
			return
		}
		if response.path != "/foo" {
			t.Errorf("response default path want %s but get %s", "/foo", response.path)
		}
		if response.Data != nil {
			t.Errorf("response default data is not nil")
		}
		ba := response.BaseAuth()
		if ba != nil {
			t.Errorf("response default baseauth is not nil")
		}
		response.SetBaseAuth("account", "pwd")
		ba = response.BaseAuth()
		if ba.Account != "account" {
			t.Errorf("response baseauth default accout is not account")
		}
		if ba.Pwd != "pwd" {
			t.Errorf("response baseauth default password is not pwd")
		}
		if !ba.Enable {
			t.Errorf("response baseauth default enable is not true")
		}
		fields := response.GetFields()
		if fields != nil {
			t.Errorf("response default fields is not nil")
		}
		f := map[string]string{
			"a": "a",
			"b": "b",
		}
		response.SetFields(f)
		fields = response.GetFields()
		if !reflect.DeepEqual(f, fields) {
			t.Error("fields is set failed")
		}
	})

	t.Run("test get file", func(t *testing.T) {
		response := NewResponse("/txt/file.txt")
		err := client.GetFile(response)
		if err != nil {
			t.Error(err.Error())
			return
		}
	})
	t.Run("test upload file", func(t *testing.T) {
		response := NewResponse("/upload")
		response.SetUploadFile("./upload.txt")
		defer response.Close()
		response.SetFields(map[string]string{"filename": "upload.txt"})
		_, err := client.Upload(response)
		if err != nil {
			t.Error(err.Error())
			return
		}
	})

	t.Run("test get", func(t *testing.T) {
		response := NewResponse("/bar")
		_, err := client.Get(response)
		if err != nil {
			t.Error(err.Error())
			return
		}
	})

	t.Run("test post", func(t *testing.T) {
		response := NewResponse("/foo")
		b, err := json.Marshal(map[string]interface{}{
			"a": "a",
			"b": "b",
		})
		if err != nil {
			t.Error(err.Error())
			return
		}
		_, err = client.Post(response, string(b))
		if err != nil {
			t.Error(err.Error())
			return
		}
	})

}
