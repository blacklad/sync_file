package utils

import (
	"os"
)

// 文件类型
type FileType string

// 事件类型
type EventType string

const (
	File FileType = "file"
	Dir  FileType = "dir"
)

// 文件信息
type FileStat struct {
	Path         string   `json:"path"`
	MD5          string   `json:"md5"`
	FileType     FileType `json:"fileType"`
	LastModified int64    `json:"lastModified"`
	Version      int64
}

func GetFileType(fileInfo os.FileInfo) FileType {
	if fileInfo.IsDir() {
		return Dir
	}
	return File
}

func FileIsExists(path string) bool {
	_, err := os.Stat(path)
	// 如果文件不存在则不处理
	if err != nil {
		return false
	}
	return true
}
