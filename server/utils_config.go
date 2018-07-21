package main

import (
	"io/ioutil"

	"github.com/tidwall/gjson"
)

var config Config

// Config class
type Config struct {
	path    string
	content string
}

// NewConfig create Config class
func NewConfig(path string) (Config, error) {
	config := Config{}
	config.path = path

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}

	str := string(bytes)
	config.content = str

	return config, nil
}

func (config *Config) get(arg ...interface{}) string {
	value := gjson.Get(config.content, arg[0].(string)).String()

	if (value == "") && len(arg) > 1 {
		value = arg[1].(string)
	}

	return value
}
