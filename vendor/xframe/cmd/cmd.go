package cmd

import (
	"errors"
	"flag"
	"fmt"
)

func GetCommand(name string) (value string, err error) {
	flag := flag.Lookup(name)
	if flag == nil {
		return "", errors.New(fmt.Sprintf("Can't find command[\"%s\"]", name))
	}
	return flag.Value.String(), nil
}

func Usage(err error) {
	if err != nil {
		fmt.Println("ERROR:", err)
	}
	flag.Usage()
}

func ParseCommand() {
	if !flag.Parsed() {
		flag.Parse()
	}
}

func DumpCommand() {
	if !flag.Parsed() {
		flag.Parse()
	}
	visitor := func(a *flag.Flag) {
		fmt.Println("option =", a.Name, " value =", a.Value)
	}
	flag.VisitAll(visitor)
}
