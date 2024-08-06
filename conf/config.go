package conf

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Ms Ms `yaml:"miaosha"`
}

type Ms struct {
	Redis    string `yaml:"redis"`
	MysqlDsn string `yaml:"mysql_dsn"`
}

var ConfigData Config

func init() {
	bytes, err := os.ReadFile("conf/config.yaml")
	if err != nil {
		log.Fatal("配置文件打开失败", err)
	}

	err = yaml.Unmarshal(bytes, &ConfigData)
	if err != nil {
		log.Fatal("配置文件解析失败", err)
	}
}
