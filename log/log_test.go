package log

import (
	"runtime"
	"testing"
)

func TestWrite(t *testing.T) {
	InitLogger(".", "test", "logtest", 10, "DEBUG", "")
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

//TODO: test case for caller
