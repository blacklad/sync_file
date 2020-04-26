package conf

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestGetConf(t *testing.T) {
	config, err := GetConf("")
	if err != nil {
		fmt.Print(err)
	}
	fmt.Println(config)
}

func TestGet(t *testing.T) {
	fmt.Println(filepath.Join("aa/bb/dd"))
}
