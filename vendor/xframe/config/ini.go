package config

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
)

type KV map[string]interface{}

var (
	configContentINI = make(map[string]KV)
)

const (
	BLOCK_LEFT  = '['
	BLOCK_RIGHT = ']'

	KV_SEPARATOR   = "="
	LINE_SEPARATOR = "\n"
)

func _parseINI(data []byte) (err error) {
	var (
		curr_block string
	)
	//init
	configContentINI[curr_block] = make(KV)
	lines := bytes.Split(data, []byte(LINE_SEPARATOR))
	for _, line := range lines {
		if len(line) == 0 {
			//space line
			continue
		}
		line = bytes.TrimSpace(line)
		if line[0] == ';' || line[0] == '#' {
			//comment
			continue
		}
		if line[0] == BLOCK_LEFT && line[len(line)-1] == BLOCK_RIGHT {
			block_name := string(line[1 : len(line)-1])
			curr_block = block_name
			if _, ok := configContentINI[block_name]; !ok {
				configContentINI[block_name] = make(KV)
			}
			continue
		}
		key := bytes.TrimSpace(bytes.Split(line, []byte(KV_SEPARATOR))[0])
		value := bytes.TrimSpace(bytes.Split(line, []byte(KV_SEPARATOR))[1])
		configContentINI[curr_block][string(key)] = string(value)
	}
	return
}

func LoadFromFileINI(filename string) (err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	return _parseINI(data)
}

func GetConfigINI(block string, key string) (value interface{}, err error) {
	if kv, ok := configContentINI[block]; ok {
		value, ok = kv[key]
		if !ok {
			err = errors.New(fmt.Sprintf("can not get key [\"%s\"]", key))
		}
	} else {
		err = errors.New(fmt.Sprintf("can not get block [\"%s\"]", block))
	}
	return
}

func DumpConfigINI() {
	for block_name, block := range configContentINI {
		fmt.Printf("[%s]\n", block_name)
		for k, v := range block {
			fmt.Printf("\t%s = %s\n", k, v.(string))
		}
	}
}
