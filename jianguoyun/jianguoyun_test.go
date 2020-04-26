package jianguoyun

import (
	"blacklad.com/sync_file/conf"
	"blacklad.com/sync_file/utils"
	"fmt"
	"testing"
)

func TestNewJianGuoYunClient(t *testing.T) {
	config, err := conf.GetConf("../conf/conf.yaml")
	utils.HandleError(err)

	path := "/"

	j := NewJianGuoYunClient(config.SyncConfig.Url, config.SyncConfig.User, config.SyncConfig.Password, path)

	fileList, _ := j.List()

	fmt.Println(len(fileList))
}

func TestJianGuoYun_DownloadFile(t *testing.T) {
	config, err := conf.GetConf("../conf/conf.yaml")
	utils.HandleError(err)
	path := "/"

	j := NewJianGuoYunClient(config.SyncConfig.Url, config.SyncConfig.User, config.SyncConfig.Password, path)

	j.DownloadFile("", "")
}
