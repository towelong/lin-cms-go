package pkg

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// 最终方案-全兼容
func GetCurrentAbPath() string {
	dir := GetCurrentAbPathByExecutable()
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	if strings.Contains(dir, tmpDir) {
		return strings.TrimRight(GetCurrentAbPathByCaller(), "/pkg")
	}
	return strings.TrimRight(dir, "/pkg")
}

// 获取当前执行文件绝对路径
func GetCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func GetCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

// 判断所给路径文件夹是否存在
func IsDirExist(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		log.Println(err)
		return false
	}
	return s.IsDir()
}

// 返回指定目录的绝对路径(不存在则创建目录)
func CreateDirAndFileForCurrentTime(fileDir string, format string) (string, error) {
	s := "/" + time.Now().Format(format)
	dir := path.Join(GetCurrentAbPath(), "/"+fileDir, s)
	if !IsDirExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	return dir, nil
}
