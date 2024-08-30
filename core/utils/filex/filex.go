package filex

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateFile(fileName string, data []byte) error {
	// 确保目录存在
	fileDir := filepath.Dir(fileName)
	if err := os.MkdirAll(fileDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %s", fileDir)
	}
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("创建文件失败: %s", fileName)
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("写入文件失败: %s", fileName)
	}
	return nil
}
