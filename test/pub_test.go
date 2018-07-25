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
		rst := make(chan interface{}, 1)
		f.Recent("测试", func(id ...interface{}) interface{} {
			fmt.Println(id)
			return id
		}, rst, i)
		fmt.Println("返回结果：", <-rst)
	}
}