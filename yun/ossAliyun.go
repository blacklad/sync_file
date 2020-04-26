package yun

import (
	"blacklad.com/sync_file/conf"
	"blacklad.com/sync_file/utils"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"path/filepath"
	"strings"
)

const OssDirSeparator = "/"

type OssAli struct {
	basePath     string
	BucketClient *oss.Bucket
}

func NewOssAli(config *conf.Config) (*OssAli, error) {
	// 创建OSSClient实例。
	client, err := oss.New(config.OssConfig.Endpoint, config.OssConfig.Key, config.OssConfig.Secret)
	if err != nil {
		utils.HandleError(err)
	}
	// 获取存储空间。
	bucketName := "blacklad"
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}

	return &OssAli{
		basePath:     config.OssBasePath,
		BucketClient: bucket,
	}, nil
}

// 上传文件到oss。
func (o *OssAli) UploadFile(filePath, localPath string) error {
	ossPath := filepath.Join(o.basePath, filePath)
	err := o.BucketClient.PutObjectFromFile(ossPath, localPath)
	return err
}

// oss里创建一个文件夹
// oss中文件夹为一个没有内容的文件且以/结尾
func (o *OssAli) CreateDir(filePath string) error {
	ossPath := filepath.Join(o.basePath, filePath) + OssDirSeparator
	err := o.BucketClient.PutObject(ossPath, strings.NewReader(""))
	return err
}

// 删除一个文件夹
func (o *OssAli) DeleteFile(filePath string) error {
	ossPath := filepath.Join(o.basePath, filePath)
	err := o.BucketClient.DeleteObject(ossPath)
	return err
}
