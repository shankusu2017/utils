package utils

import (
	"os"
	"path/filepath"
)

// IsFileDirExist 文件或者目录存在吗？
func IsFileDirExist(fileName string) bool {
	_, error := os.Stat(fileName)
	// check if error is "file not exists"
	if os.IsNotExist(error) {
		//fmt.Printf("%v file does not exist\n", fileName)
		return false
	}

	return true
}

// MakeFile 文件不存在则创建，存在则直接返回
func MakeFile(allPath string) error {
	if IsFileDirExist(allPath) {
		return nil
	}

	// 创建前置目录
	preDir := filepath.Dir(allPath)
	err := os.MkdirAll(preDir, 0755)
	if err != nil {
		return err
	}

	// 创建文件
	fd, err := os.OpenFile(allPath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fd.Close()

	return nil
}

// OpenFile 可读写的方式打开(若不存在则创建)文件
// 写时自动追加在末尾
func OpenFile(allPath string) (*os.File, error) {
	err := MakeFile(allPath)
	if err != nil {
		return nil, err
	}

	fd, err := os.OpenFile(allPath, os.O_RDWR|os.O_APPEND, 0644)
	return fd, err
}
