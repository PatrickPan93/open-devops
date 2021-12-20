package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	RPCServerAddr string `yaml:"rpc_server_addr"`
}

// Load 根据LoadFile读取配置文件后的字符串解析yaml为配置结构体
func Load(bs []byte) (*Config, error) {
	cfg := &Config{}
	err := yaml.Unmarshal(bs, cfg)
	if err != nil {
		return nil, errors.Wrap(err, "Load: Loading file filed")
	}
	return cfg, nil
}

// LoadFile 根据conf路径读取内容
func LoadFile(filename string) (*Config, error) {

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "LoadFile: Error while reading file via ReadFile")
	}

	cfg, err := Load(bytes)

	if err != nil {
		return nil, errors.Wrap(err, "LoadFile: Error while reader reading bytes")
	}
	return cfg, nil
}
