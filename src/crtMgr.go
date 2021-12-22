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
	crtPool []*crt
	buckets []*readyCrtsBuckets
	Wait    *sync.WaitGroup
}

type readyCrtsBuckets struct {
	readyCrts      map[int]*crt
	cap            int
	readyCrtsMutex *sync.RWMutex
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

// CreateCrts ...
func (this *crtMgr) CreateCrts(bucketNums int, numsEachBucket int) {
	rand.Seed(time.Now().UnixNano())
	this.buckets = make([]*readyCrtsBuckets, bucketNums)
	this.crtPool = make([]*crt, bucketNums*numsEachBucket)

	for bucketId := 0; bucketId < bucketNums; bucketId++ {
		bucket := new(readyCrtsBuckets)
		bucket.cap = numsEachBucket
		bucket.readyCrts = make(map[int]*crt)
		bucket.readyCrtsMutex = new(sync.RWMutex)
		this.buckets[bucketId] = bucket
		for id := bucketId * numsEachBucket; id < (bucketId+1)*numsEachBucket; id++ {
			fooCrt := newCrt(id, bucketId)
			this.crtPool[id] = fooCrt
		}

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

	for id, bucket := range this.buckets {
		for bucket.cap != len(bucket.readyCrts) {
			time.Sleep(100 * time.Millisecond)
			fmt.Println("bucket[", id, "]:", "total: ", bucket.cap, "readyCrts:", len(bucket.readyCrts))
		}
	}

	killTask := new(task)
	killTask.release = true

	for _, crt := range this.crtPool {
		crt.task = killTask
		crt.wakeup()
	}
}

// pop a ready coroutine
func (this *crtMgr) popReadyCrt() *crt {
	var crtIdx int

	randBucketId := rand.Intn(len(this.buckets))

	bucket := this.buckets[randBucketId]

	bucket.readyCrtsMutex.Lock()

	for len(bucket.readyCrts) <= 0 {
		bucket.readyCrtsMutex.Unlock()
		time.Sleep(100 * time.Millisecond)
		bucket.readyCrtsMutex.Lock()
	}

	for key, _ := range bucket.readyCrts {
		crtIdx = key
		break
	}

	delete(bucket.readyCrts, crtIdx)
	bucket.readyCrtsMutex.Unlock()

	crt := this.crtPool[crtIdx]
	return crt
}
