package jianguoyun

import (
	"blacklad.com/sync_file/utils"
	"io"
	"os"
	"path/filepath"
	"studio-b12/gowebdav"
)

var FilePermMode = os.FileMode(0777) // Default file permission

type JianGuoYun struct {
	user     string
	password string
	url      string

	path   string
	client *gowebdav.Client
}

func NewJianGuoYunClient(url, user, password, path string) *JianGuoYun {
	client := gowebdav.NewClient(url, user, password)

	j := &JianGuoYun{
		url:      url,
		user:     user,
		password: password,
		path:     path,
		client:   client,
	}

	return j
}

func (j JianGuoYun) List() ([]*utils.FileStat, error) {
	var fileList = make([]*utils.FileStat, 0)
	j.list(&fileList, j.path)
	return fileList, nil
}

func (j JianGuoYun) list(fileList *[]*utils.FileStat, path string) {
	files, _ := j.client.ReadDir(path)
	for _, file := range files {
		filePath := filepath.Join(path, file.Name())

		if path == filePath {
			continue
		}

		if file.IsDir() {
			f := &utils.FileStat{
				Path:         filePath,
				MD5:          file.ETag(),
				FileType:     utils.Dir,
				LastModified: file.ModTime().Unix(),
			}
			*fileList = append(*fileList, f)
			j.list(fileList, filePath)
		} else {
			f := &utils.FileStat{
				Path:         filePath,
				MD5:          file.ETag(),
				FileType:     utils.File,
				LastModified: file.ModTime().Unix(),
			}
			*fileList = append(*fileList, f)
		}
	}
}

func (j JianGuoYun) DownloadFile(jgyPath, tmpPath string) (bool, error) {
	tmpDir := filepath.Dir(tmpPath)
	if !utils.FileIsExists(tmpDir) {
		err := os.MkdirAll(tmpDir, FilePermMode)
		if err != nil {
			return false, err
		}
	}

	data, err := j.client.ReadStream(jgyPath)
	if _, ok := err.(*os.PathError); ok {
		// 文件不存在
		return false, nil
	}

	if err != nil {
		return false, err
	}

	// If the local file does not exist, create a new one. If it exists, overwrite it.
	fd, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, FilePermMode)
	if err != nil {
		return false, err
	}

	// Copy the data to the local file path.
	_, err = io.Copy(fd, data)
	if err != nil {
		return false, err
	}

	defer fd.Close()

	return true, nil
}
