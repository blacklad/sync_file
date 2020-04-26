package storage

import (
	"blacklad.com/sync_file/conf"
	"blacklad.com/sync_file/utils"
	"fmt"
	"testing"
)

func TestNewDbFile(t *testing.T) {
	config, err := conf.GetConf("../conf/conf.yaml")
	utils.HandleError(err)

	d, _ := NewDbFile(config.LocalBasePath)
	d.createTable()

	f := &DbFileStat{
		FileStat: utils.FileStat{
			Path:         "b",
			MD5:          "b",
			FileType:     "c",
			LastModified: 1234,
			Version:      100,
		},
	}

	d.Add(f)

	//d.DeleteById(2)
	//d.UpdateById(1, f)

	d.List()
	for i := 0; i < len(d.fileList); i++ {
		fmt.Println(d.fileList[i])
	}

}

func TestDbFile_GetMaxVersion(t *testing.T) {
	config, err := conf.GetConf("../conf/conf.yaml")
	utils.HandleError(err)

	d, _ := NewDbFile(config.LocalBasePath)

	version, err := d.GetMaxVersion()
	utils.HandleError(err)
	fmt.Println(version)

}

func TestCreateDB(t *testing.T) {
	config, err := conf.GetConf("../conf/conf.yaml")
	utils.HandleError(err)
	d, _ := NewDbFile(config.LocalBasePath)

	err = d.createTable()
	if err != nil {
		fmt.Println(err)
	}

}
