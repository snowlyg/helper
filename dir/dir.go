package dir

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v2"
)

// GetCurrentAbPath 项目根目录绝对路径
func GetCurrentAbPath() string {
	dir := GetCurrentAbPathByExecutable()
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	if strings.Contains(dir, tmpDir) {
		return GetCurrentAbPathByCaller()
	}
	if strings.Contains(dir, "/tmp") {
		return filepath.Dir(dir)
	}

	return dir
}

// GetCurrentAbPathByExecutable 当前执行文件目录
func GetCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// GetCurrentAbPathByCaller 当前方法执行目录
func GetCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

// GetCurrentFuncNameByCaller 当前方法执行函数名
func GetCurrentFuncNameByCaller() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	return f.Name()
}

// RealPath 基于构件执行文件的绝对文件路径
func RealPath(fp string) (string, error) {
	if path.IsAbs(fp) {
		return fp, nil
	}
	wd, err := os.Getwd()
	return path.Join(wd, fp), err
}

// SelfPath 完整的执行文件绝对路径
func SelfPath() string {
	path, _ := filepath.Abs(os.Args[0])
	return path
}

// SelfDir 执行文件目录完整路径
func SelfDir() string {
	return filepath.Dir(SelfPath())
}

// Basename 从路径中提取文件名
func Basename(fp string) string {
	return path.Base(fp)
}

// Dir 从路径中获取目录路径
func Dir(fp string) string {
	return path.Dir(fp)
}

//InsureDir 新建不存在的文件夹
func InsureDir(fp string) error {
	if IsExist(fp) {
		return nil
	}
	return os.MkdirAll(fp, os.ModePerm)
}

// IsExist 检测文件或者目录是否存在
// 不存在的时候将返回 fasle
func IsExist(fp string) bool {
	_, err := os.Stat(fp)
	return err == nil || os.IsExist(err)
}

// IsFile checks whether the path is a file,
// it returns false when it's a directory or does not exist.
func IsFile(fp string) bool {
	f, e := os.Stat(fp)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

func ReadBytes(cpath string) ([]byte, error) {
	if !IsExist(cpath) {
		return nil, fmt.Errorf("%s not exists", cpath)
	}

	if !IsFile(cpath) {
		return nil, fmt.Errorf("%s not file", cpath)
	}

	return ioutil.ReadFile(cpath)
}

func ReadString(cpath string) (string, error) {
	bs, err := ReadBytes(cpath)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}

func ReadStringTrim(cpath string) (string, error) {
	out, err := ReadString(cpath)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out), nil
}

func ReadYaml(cpath string, cptr interface{}) error {
	bs, err := ReadBytes(cpath)
	if err != nil {
		return fmt.Errorf("cannot read %s: %s", cpath, err.Error())
	}

	err = yaml.Unmarshal(bs, cptr)
	if err != nil {
		return fmt.Errorf("cannot parse %s: %s", cpath, err.Error())
	}

	return nil
}

func ReadJson(cpath string, cptr interface{}) error {
	os.MkdirAll(path.Dir(cpath), os.ModePerm)
	bs, err := ReadBytes(cpath)
	if err != nil {
		return fmt.Errorf("cannot read %s: %s", cpath, err.Error())
	}

	err = json.Unmarshal(bs, cptr)
	if err != nil {
		return fmt.Errorf("cannot parse %s: %s", cpath, err.Error())
	}

	return nil
}

func WriteBytes(filePath string, b []byte) (int, error) {
	os.MkdirAll(path.Dir(filePath), os.ModePerm)
	fw, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer fw.Close()
	return fw.Write(b)
}

func WriteString(filePath string, s string) (int, error) {
	return WriteBytes(filePath, []byte(s))
}

func MD5(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	h := md5.New()

	_, err = io.Copy(h, r)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func Md5Byte(p []byte) (string, error) {
	h := md5.New()
	_, err := h.Write(p)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func OpenLogFile(fp string) (*os.File, error) {
	os.MkdirAll(path.Dir(fp), os.ModePerm)
	return os.OpenFile(fp, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

// 添加文本
func AppendFile(filePath string, b []byte) error {
	os.MkdirAll(path.Dir(filePath), os.ModePerm)
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(string(b) + "\r\n\r\n")

	return nil
}
