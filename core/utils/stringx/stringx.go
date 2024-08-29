package stringx

import (
	"github.com/golang-module/carbon/v2"
	"log"
	"regexp"
	"strconv"
	"strings"
)

// ByteSize 体积字符串转字节数
func ByteSize(str string) int64 {
	if str == "" {
		return 0
	}
	// 去除非法字符，并转换为大写
	str = strings.TrimSpace(str)
	str = strings.ReplaceAll(str, ",", "")
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, "\n", "")
	str = strings.ToUpper(str)
	// 匹配大小单位（如 KB、MB、GB 等）
	re := regexp.MustCompile(`[KMGTPI]*B?`)
	sizeStr := re.ReplaceAllString(str, "")
	// 尝试转换为浮点数
	size, err := strconv.ParseFloat(sizeStr, 64)
	if err != nil {
		log.Printf("字符串转 float64 异常: %v", err)
		return 0
	}
	// 根据单位乘以相应的倍数
	switch {
	case strings.Contains(str, "PB") || strings.Contains(str, "PIB") || strings.Contains(str, "P"):
		size *= 1024 * 1024 * 1024 * 1024 * 1024
	case strings.Contains(str, "TB") || strings.Contains(str, "TIB") || strings.Contains(str, "T"):
		size *= 1024 * 1024 * 1024 * 1024
	case strings.Contains(str, "GB") || strings.Contains(str, "GIB") || strings.Contains(str, "G"):
		size *= 1024 * 1024 * 1024
	case strings.Contains(str, "MB") || strings.Contains(str, "MIB") || strings.Contains(str, "M"):
		size *= 1024 * 1024
	case strings.Contains(str, "KB") || strings.Contains(str, "KIB") || strings.Contains(str, "K"):
		size *= 1024
	}
	// 四舍五入并返回整数值
	return int64(size + 0.5)
}

// ParseInt 将字符串转换为整数
func ParseInt(str string) int {
	if len(str) == 0 {
		return 0
	}
	str = strings.ReplaceAll(str, ",", "")
	value, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("字符串转 int 异常: %v", err)
		return 0
	}
	return value
}

func ParseInt64(str string) int64 {
	if len(str) == 0 {
		return 0
	}
	str = strings.ReplaceAll(str, ",", "")
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Printf("字符串转 int64 异常: %v", err)
		return 0
	}
	return i
}

func ParseFloat64(str string) float64 {
	if len(str) == 0 {
		return 0
	}
	str = strings.ReplaceAll(str, ",", "")
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Printf("字符串转 float64 异常: %v", err)
		return 0
	}
	return f
}

func ParseBool(str string) bool {
	if len(str) == 0 {
		return false
	}
	b, err := strconv.ParseBool(str)
	if err != nil {
		log.Printf("字符串转 bool 异常: %v", err)
		return false
	}
	return b
}

// TimeStamp 解析日期字符串并返回时间戳
func TimeStamp(date string) int64 {
	if len(date) == 0 {
		return 0
	}
	return carbon.Parse(date).Timestamp()
}
