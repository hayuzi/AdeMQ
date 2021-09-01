package conf

import (
	"flag"
	"fmt"
	"github.com/AdeMQ/server/service"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

var (
	Conf     = &Config{}
	confPath string
)

type Config struct {
	Server *service.Config
	Logger *Logger
}

type Server struct {
	Address string `yaml:"address" json:"address"`
	Port    string `yaml:"port" json:"port"`
}

type Logger struct {
	Stdout bool `yaml:"stdout" json:"stdout"`
	File   *File
}

type File struct {
	Dir  string `yaml:"dir" json:"dir"`
	Type string `yaml:"type" json:"type"`
}

func init() {
	flag.StringVar(&confPath, "conf", "", "conf values")
}

func Init() (err error) {
	fmt.Println(filepath.Abs("./config.yaml"))
	var (
		yamlFile string
	)
	if confPath != "" {
		yamlFile, err = filepath.Abs(confPath)
	} else {
		yamlFile, err = filepath.Abs("./config.yaml")
	}
	if err != nil {
		return
	}
	yamlRead, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(yamlRead, Conf)
	if err != nil {
		return
	}
	go load()
	return
}

// 动态加载配置
func load() {

}
