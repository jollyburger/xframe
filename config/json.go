package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

var (
	configContentByte []byte
	configContentJson map[string]interface{}
)

type block struct {
	data interface{}
}

func getConfig(data map[string]interface{}) (*block, error) {
	if data == nil {
		return nil, errors.New("get config fail, config content is nil")
	}
	return &block{
		data: data,
	}, nil
}

func (b *block) getValue(key string) *block {
	m := b.getData()
	if v, ok := m[key]; ok {
		b.data = v
		return b
	}
	return nil
}

func (b *block) getData() map[string]interface{} {
	if m, ok := (b.data).(map[string]interface{}); ok {
		return m
	}
	return nil
}

// Load config from file
func LoadConfigFromFile(filename string) (err error) {
	configContentByte, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	return json.Unmarshal(configContentByte, &configContentJson)
}

// delim: ".", get kv from config
func GetConfigByKey(keys string) (value interface{}, err error) {
	key_list := strings.Split(keys, ".")
	block, err := getConfig(configContentJson)
	if err != nil {
		return nil, err
	}
	for _, key := range key_list {
		block = block.getValue(key)
		if block == nil {
			return nil, errors.New(fmt.Sprintf("can not get[\"%s\"]'s value", string(keys)))
		}
	}
	return block.data, nil
}

func DumpConfigContent() {
	var pjson bytes.Buffer
	json.Indent(&pjson, configContentByte, "", "\t")
	fmt.Println(string(pjson.Bytes()))
}

//v0.3: read config json file with customized struct
func LoadConfigFromFileV2(config interface{}, filepath string) (err error) {
	confBuf, err := ioutil.ReadFile(filepath)
	if err != nil {
		return
	}
	err = json.Unmarshal(confBuf, config)
	return
}
