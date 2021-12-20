package GCRTPool

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var mgr *crtMgr

var single sync.Once

type crtMgr struct {
	crtPool        []*crt
	readyCrts      map[int]*crt
	readyCrtsMutex sync.RWMutex
	Wait           *sync.WaitGroup
}

// GetMgr get the CrtMgr; CrtMgr is a singleton
func GetMgr() *crtMgr {
	single.Do(
		func() {
			mgr = &crtMgr{
				Wait: new(sync.WaitGroup),
			}
		})
	return mgr
}

// create coroutine ...
func (this *crtMgr) CreateCrts(nums int) {
	rand.Seed(time.Now().UnixNano())
	this.readyCrts = make(map[int]*crt)
	for i := 0; i < nums; i++ {
		//TODO: create coroutine
		fooCrt := newCrt()
		fmt.Println("creatcrt", i, "success")
		this.crtPool = append(this.crtPool, fooCrt)
		//this.crtPool[i] = fooCrt
		// this.readyCrts[i] = fooCrt
	}

	for nums != len(this.readyCrts) {
		time.Sleep(200 * time.Millisecond)
		fmt.Println("total:", nums, "readyCrts:", len(this.readyCrts))
	}

}

// Process one task ...
func (this *crtMgr) Process(t *task) bool {
	foo := this.popReadyCrt()
	//fmt.Println("task send to foo", foo.id)
	foo.cond.L.Lock()
	foo.task = t
	foo.stat = false
	foo.cond.L.Unlock()
	foo.wakeup()
	//fmt.Println("signal ", foo.id, " ", time.Now())
	return true
}

func (this *crtMgr) KillAllCoroutine(params []interface{}) {
	total := len(this.crtPool)
	for total != len(this.readyCrts) {
		time.Sleep(100 * time.Millisecond)
		//fmt.Println("total:", total, "readyCrts:", len(this.readyCrts))
	}

	for i := 0; i < total; i++ {
		foo := new(task)
		foo.release = true
		this.Process(foo)
	}
}

// pop a ready coroutine
func (this *crtMgr) popReadyCrt() *crt {
	var fooIdx int
	this.readyCrtsMutex.Lock()

	for len(this.readyCrts) <= 0 {
		this.readyCrtsMutex.Unlock()
		time.Sleep(100 * time.Millisecond)
		this.readyCrtsMutex.Lock()
	}

	for key, _ := range this.readyCrts {
		fooIdx = key
		break
	}

	delete(this.readyCrts, fooIdx)
	this.readyCrtsMutex.Unlock()

	foo := this.crtPool[fooIdx]
	return foo
}
