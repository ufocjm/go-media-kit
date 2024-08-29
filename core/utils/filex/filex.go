package filex

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateFile(dir string, fileName string, data []byte) error {
	f := filepath.Join(dir, fileName)
	// 确保目录存在
	fileDir := filepath.Dir(f)
	if err := os.MkdirAll(fileDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %s", fileDir)
	}
	file, err := os.Create(f)
	if err != nil {
		return fmt.Errorf("创建文件失败: %s", f)
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("写入文件失败: %s", f)
	}
	return nil
}
