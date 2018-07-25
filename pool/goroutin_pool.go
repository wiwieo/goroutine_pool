package pool

import (
	"sync"
	"time"
)

const (
	FACTORY_SCALE = 10000
)

// 任务
type job struct {
	desc   string               // 工作描述
	do     func(...interface{}) // 工作
	params []interface{}        // 材料
}

type human struct {
	id   int
	name string
	busy bool
}

// 工人
type worker struct {
	human
	jobs         chan job
	wantToTravel chan bool
}

// 人事（管理工人，招聘及离职）
type hr struct {
	work    chan worker // 用于通知工人开工
	workers []*worker   // 用于确定当前有哪些工人可以工作
	m       sync.Mutex
}

func (h *hr) leave() {
	for _, w := range h.workers {
		go func() {
			<-w.wantToTravel
		}()
	}
}

// 一个工人同一时间只能干一件事，所以是阻塞的
func (w *worker) work() {
	isBreak := false
	for {
		select {
		// 等待安排工作
		case j := <-w.jobs: // 来工作立即开工
			j.do(j.params)
		case <-time.After(10 * time.Second): // 10秒没工作后，工人自动离职
			w.wantToTravel <- true
			isBreak = true
		}
		if isBreak {
			break
		}
	}
}

// 工厂（对外提供服务）
type factory struct {
	h hr
}

// 开业
func Open() *factory {
	// 建一个工厂
	f := &factory{
		h: hr{
			work: make(chan worker, FACTORY_SCALE),
		},
	}

	// 工人随时准备开工
	go func(f *factory) {
		for {
			w := <-f.h.work
			go w.work()
		}
	}(f)
	return f
}

// 接单
func (f *factory) Recent(name string, do func(...interface{}), params ...interface{}) {
	go func() {
		w := f.h.employee()
		w.busy = true
		f.h.work <- *w
		// 告诉工人来工作了
		w.jobs <- job{
			do:     do,
			params: params,
		}
		f.h.workers[w.id].busy = false
	}()
}

// 找空闲的工人
func (h *hr) findWorker() *worker {
	for _, w := range h.workers {
		if w != nil && !w.busy {
			return w
		}
	}
	return nil
}

// 招工(先从内容招，如果没有，再重新创建一个)
func (h *hr) employee() *worker {
	var w *worker
	// 如果满员且都在工作，则等待，直到有空闲的工人为止
	for {
		w = h.findWorker()
		if w != nil {
			return w
		} else {
			if len(h.workers) < FACTORY_SCALE {
				h.m.Lock()
				w = &worker{jobs: make(chan job, 1), human: human{id: len(h.workers)}}
				h.workers = append(h.workers, w)
				h.m.Unlock()
				return w
			}
			// 此处必须稍作停顿，否则容易出现循环不止，消耗系统资源
			time.Sleep(time.Millisecond)
		}
	}
}
