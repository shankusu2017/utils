package utils

// IsFileDirExist 文件或者目录存在吗？
func IsFileDirExist(fileName string) bool {
	_ , error := os.Stat(fileName)
	// check if error is "file not exists"
	if os.IsNotExist(error) {
		//fmt.Printf("%v file does not exist\n", fileName)
		return false
	}

	return true
}

