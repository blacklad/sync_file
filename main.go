package main

import (
	"blacklad.com/sync_file/conf"
	"blacklad.com/sync_file/jianguoyun"
	"blacklad.com/sync_file/storage"
	"blacklad.com/sync_file/sync"
	"blacklad.com/sync_file/utils"
	"blacklad.com/sync_file/yun"
	"fmt"
	"time"
)

func main() {
	config, err := conf.GetConf("/Users/blacklad/go/src/blacklad.com/webdav/conf/conf.yaml")
	utils.HandleError(err)

	syncClient := jianguoyun.NewJianGuoYunClient(config.SyncConfig.Url, config.SyncConfig.User, config.SyncConfig.Password, config.SyncBasePath)

	dbClient, err := storage.NewDbFile(config.DbPath)
	utils.HandleError(err)

	yunClient, err := yun.NewOssAli(config)
	utils.HandleError(err)

	syncTime := time.Second * 100
	timeSyncFile, err := sync.NewTimeSyncFile(syncTime, config.LocalBasePath, syncClient, dbClient, yunClient)
	utils.HandleError(err)

	timeSyncFile.Start()

	errCh := make(chan error, 1)
	select {
	case err := <-errCh:
		fmt.Println("error", err)
	}
}
