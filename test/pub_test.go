package test

import (
	"fmt"
	"log"
	"os"
	"pool"
	"pool_new"
	"runtime"
	"runtime/pprof"
	"testing"
)

func Benchmark_Pool_Result(b *testing.B) {
	b.SetParallelism(pool.FACTORY_SCALE)
	f := pool.Open()
	for i := 0; i < b.N; i++ {
		rst := make(chan interface{}, 1)
		f.Recent("测试", add, rst, i, i + 1)
		fmt.Println("返回结果：", <-rst)
	}
}

func Benchmark_Pool(b *testing.B) {
	b.SetParallelism(pool.FACTORY_SCALE)
	f := pool.Open()

	for i := 0; i < b.N; i++ {
		f.Recent("测试", add, nil, i, i + 1)
	}
	//Profile()
}

func Benchmark_Pool_New_Result(b *testing.B) {
	b.SetParallelism(pool.FACTORY_SCALE)
	f := pool_new.NewFactoryAndRun()

	for i := 0; i < b.N; i++ {
		rst := make(chan interface{}, 1)
		f.Recent("测试", add, rst, i, i + 1)
		fmt.Println(i, "返回结果：", <-rst)
	}
	//Profile()
}

func add(i ...interface{}) interface{} {
	vi, _ := i[0].(int)
	vj, _ := i[1].(int)
	return vi + vj
}

func Benchmark_Pool_New(b *testing.B) {
	b.SetParallelism(pool.FACTORY_SCALE)
	f := pool_new.NewFactoryAndRun()

	for i := 0; i < b.N; i++ {
		f.Recent("测试", add, nil, i, i+1)
	}
	//Profile()
}

func Profile() {
	f, err := os.Create("mem_profile.prof")
	if err != nil {
		log.Fatal("could not create memory profile: ", err)
	}
	runtime.GC() // get up-to-date statistics
	if err := pprof.WriteHeapProfile(f); err != nil {
		log.Fatal("could not write memory profile: ", err)
	}
	f.Close()
}

func Benchmark_Print(b *testing.B) {
	b.SetParallelism(pool.FACTORY_SCALE)
	for i := 0; i < b.N; i++ {
		go fmt.Println(i)
	}
}
