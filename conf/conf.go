package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const default_config_yaml = "./conf.yaml"

type OssConfig struct {
	Endpoint string `yaml:"endpoint"`
	Key      string `yaml:"key"`
	Secret   string `yaml:"secret"`
}

type SyncConfig struct {
	Url      string `yaml:"url"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type Config struct {
	LocalBasePath string     `yaml:"localBasePath"`
	OssBasePath   string     `yaml:"ossBasePath"`
	SyncBasePath  string     `yaml:"syncBasePath"`
	DbPath        string     `yaml:"dbPath"`
	OssConfig     OssConfig  `yaml:"oss"`
	SyncConfig    SyncConfig `yaml:"sync"`
}

func GetConf(path string) (*Config, error) {
	if len(path) == 0 {
		path = default_config_yaml
	}
	return getConf(path)
}

//并转换成conf对象
func getConf(path string) (*Config, error) {
	conf := &Config{}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}
