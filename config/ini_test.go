package config

import (
	"fmt"
	"testing"
)

func Test_Ini(t *testing.T) {
	err := LoadFromFileINI("./test.ini")
	if err != nil {
		t.Error(err)
	}
	//dump all config ini file
	DumpConfigINI()
	//get block value
	v, _ := GetConfigINI("test", "test1")
	fmt.Println(v)
	//get non-block value
	v, _ = GetConfigINI("", "solo")
	fmt.Println(v)
}
