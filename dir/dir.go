package dir

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

// GetCurrentAbPath 项目根目录绝对路径
func GetCurrentAbPath() string {
	dir := GetCurrentAbPathByExecutable()
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	if strings.Contains(dir, tmpDir) {
		return GetCurrentAbPathByCaller()
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
