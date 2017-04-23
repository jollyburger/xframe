package trace

import (
	"fmt"
	"os"
	"runtime/trace"
)

//trace for go tool trace
func InitTraceTool() {
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()
	fmt.Println("start trace tool")
}
