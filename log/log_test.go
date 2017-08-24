package log

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestWrite(t *testing.T) {
	InitLogger(".", "test", "testlog", 10, "DEBUG", "")
	INFO("test")
}

func BenchmarkWrite(b *testing.B) {
	InitLogger(".", "test", "logtest", 100, "DEBUG", "")

	for i := 0; i < b.N; i++ {
		mem := new(runtime.MemStats)
		runtime.ReadMemStats(mem)
		INFO("every log mem alloced: ", mem.Alloc)
		INFO(i)
		runtime.ReadMemStats(mem)
		INFO("after log mem alloced: ", mem.Alloc)
	}
	runtime.GC()
	mem := new(runtime.MemStats)
	runtime.ReadMemStats(mem)
	INFO("after gc mem alloced: ", mem.Alloc)
}

func TestInitLogger(t *testing.T) {
	InitLogger("", "", "", 0, "", "stdout")
	for i := 0; i < 1000; i++ {
		INFO("test", i)
	}
}

func Test_WriteBackup(t *testing.T) {
	InitLogger(".", "test", "", 2500, "DEBUG", "")
	SetBackup(5)
	for {
		INFO(fmt.Sprintf("%d: test", time.Now().Unix()))
		time.Sleep(1 * time.Second)
	}
}
