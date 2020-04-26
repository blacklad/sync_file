package sync

import (
	"blacklad.com/sync_file/jianguoyun"
	"blacklad.com/sync_file/storage"
	"blacklad.com/sync_file/utils"
	"blacklad.com/sync_file/yun"
	"fmt"
	"path/filepath"
	"time"
)

type TimeSyncFile struct {
	syncTime time.Duration

	tmpPath string
	version int64

	syncClient *jianguoyun.JianGuoYun
	localDb    *storage.DbFile
	yunClient  *yun.OssAli
}

func NewTimeSyncFile(t time.Duration, tmpPath string, syncClient *jianguoyun.JianGuoYun, localDb *storage.DbFile, yunClient *yun.OssAli) (*TimeSyncFile, error) {
	version, err := localDb.GetMaxVersion()
	if err != nil {
		return nil, err
	}
	version++
	return &TimeSyncFile{
		syncTime:   t,
		tmpPath:    tmpPath,
		version:    version,
		syncClient: syncClient,
		localDb:    localDb,
		yunClient:  yunClient,
	}, nil
}

// 定时遍历文件
// 与内存中的修改时间以及tag对比
// 添加同步版本
// 如果没有则下载 上传到wos    添加db
// 如果不同则下载同时上传到wos  修改db
// 如果db版本小于最新版本      删除oss文件
func (s *TimeSyncFile) Start() {

	//t := time.Tick(s.syncTime)
	//for {
	//	select {
	//	case <-t:
	//
	//	}
	//}

	s.sync()
	s.findDeleteFile()
}

func (s *TimeSyncFile) sync() {
	fileList, err := s.syncClient.List()
	utils.HandleError(err)

	for i := 0; i < len(fileList); i++ {
		jgyFileStat := fileList[i]

		dbFileStat, err := s.localDb.GetByPath(fileList[i].Path)
		utils.HandleError(err)

		if dbFileStat != nil {
			// 判断修改时间是否大于db修改时间 如果是文件判断md5是否一致
			if jgyFileStat.LastModified > dbFileStat.LastModified && jgyFileStat.MD5 != dbFileStat.MD5 {
				fmt.Println("更新文件: ", dbFileStat.Path)
				// 下载文件上传到oss
				// 更新数据库
				dbFileStat = &storage.DbFileStat{
					Id: dbFileStat.Id,
					FileStat: utils.FileStat{
						Path:         jgyFileStat.Path,
						MD5:          jgyFileStat.MD5,
						FileType:     jgyFileStat.FileType,
						LastModified: jgyFileStat.LastModified,
						Version:      s.version,
					},
				}

				if dbFileStat.FileType == utils.Dir {
					err = s.yunClient.CreateDir(dbFileStat.Path)
					utils.HandleError(err)
				} else {
					localPath := filepath.Join(s.tmpPath, jgyFileStat.Path)
					res, err := s.syncClient.DownloadFile(jgyFileStat.Path, localPath)
					utils.HandleError(err)
					if !res {
						continue
					}

					err = s.yunClient.UploadFile(dbFileStat.Path, localPath)

					utils.HandleError(err)
				}

				raws, err := s.localDb.UpdateById(dbFileStat.Id, dbFileStat)
				utils.HandleError(err)

				if raws == 0 {
					fmt.Println("数据库更新失败,请检查: ", dbFileStat.Path)
				}
			} else {
				//更新本地version
				dbFileStat.Version = s.version
				raws, err := s.localDb.UpdateById(dbFileStat.Id, dbFileStat)
				utils.HandleError(err)

				if raws == 0 {
					fmt.Println("数据库删除失败,请检查: ", dbFileStat.Path)
				}
			}
		} else {
			// 下载文件上传到oss
			dbFileStat = &storage.DbFileStat{
				FileStat: utils.FileStat{
					Path:         jgyFileStat.Path,
					MD5:          jgyFileStat.MD5,
					FileType:     jgyFileStat.FileType,
					LastModified: jgyFileStat.LastModified,
					Version:      s.version,
				},
			}

			fmt.Println("添加文件: ", dbFileStat.Path)

			if dbFileStat.FileType == utils.Dir {
				err = s.yunClient.CreateDir(dbFileStat.Path)
				utils.HandleError(err)
			} else {
				localPath := filepath.Join(s.tmpPath, jgyFileStat.Path)
				res, err := s.syncClient.DownloadFile(jgyFileStat.Path, localPath)
				utils.HandleError(err)
				if !res {
					continue
				}

				err = s.yunClient.UploadFile(dbFileStat.Path, localPath)
				utils.HandleError(err)
			}
			raws, err := s.localDb.Add(dbFileStat)
			utils.HandleError(err)

			if raws == 0 {
				fmt.Println("数据库添加失败,请检查: ", dbFileStat.Path)
			}

		}
	}
}

func (s *TimeSyncFile) findDeleteFile() {
	fileList, err := s.localDb.GetByHistoryVersion(s.version)
	utils.HandleError(err)

	for i := 0; i < len(fileList); i++ {
		dbFileStat := fileList[i]
		fmt.Println("删除文件", &dbFileStat)

		err := s.yunClient.DeleteFile(dbFileStat.Path)
		utils.HandleError(err)

		_, err = s.localDb.DeleteById(dbFileStat.Id)
		utils.HandleError(err)
	}

}
