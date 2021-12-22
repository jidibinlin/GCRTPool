package GCRTPool

import (
	"fmt"
	"sync"
)

type crt struct {
	id       int
	bucketId int
	task     *task
	stat     bool //就绪状态 true 为就绪 false 为阻塞状态
	mutex    sync.Mutex
	cond     *GCRTCond
}

// NewCrt create and return a crt
func newCrt(id int, bucketId int) *crt {
	foo := new(crt)
	foo.id = id
	foo.bucketId = bucketId
	foo.cond = newCond(&foo.mutex)
	foo.stat = false
	GetMgr().Wait.Add(1)
	go foo.startCrt()
	return foo
}

func (this *crt) startCrt() {
	fmt.Println("startCrt")
	// first := 1
	for true {
		// if first == 1 {
		// 	this.stat = true
		// 	first++
		// }
		this.cond.L.Lock()
		this.wait()
		if this.task != nil {
			if this.task.release {
				//fmt.Println("kill one coroutine")
				break
			}
			//fmt.Println("crt", this.id, "process")
			this.task.function((this.task.params)...)
			this.task = nil
			this.cond.L.Unlock()
		}
	}
	GetMgr().Wait.Done()
}

func (this *crt) wait() {
	this.cond.wait(this)
}

func (this *crt) wakeup() {
	this.cond.signal()
}
