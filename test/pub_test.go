package test

import (
	"fmt"
	"testing"
	"pool"
)

func Benchmark_Pool(b *testing.B) {
	b.SetParallelism(pool.FACTORY_SCALE)
	f := pool.Open()
	for i := 0; i < b.N; i++ {
		f.Recent("测试", func(id ...interface{}) {
			fmt.Println(id)
		}, i)
	}
}