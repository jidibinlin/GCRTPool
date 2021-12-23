package GCRTPool

import (
	"sync"
	"testing"
	"time"

	GCRTPool "github.com/jidibinlin/GCRTPool/src"
)

const (
	RunTimes = 100000
	threads  = 1000
)

// func test() {
// 	for i := 0; i < RunTimes; i++ {
// 		GCRTPool.CoRun(testOne)
// 	}
// }

func testBenchOne() {
	time.Sleep(time.Duration(10) * time.Millisecond)
}

// BenchmarkGCRTPool ...
func BenchmarkGCRTPool(b *testing.B) {
	mgr := GCRTPool.NewMgr()
	mgr.CreateCrts(2000, 100)

	var wg sync.WaitGroup

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(RunTimes*threads + threads)
		for j := 0; j < threads; j++ {
			mgr.CoRun(func() {
				for i := 0; i < RunTimes; i++ {
					mgr.CoRun(func() {
						testBenchOne()
						wg.Done()
					})
				}
				wg.Done()
			})
		}
		wg.Wait()
	}
	b.StopTimer()
	mgr.KillAllCoroutine(nil)
	mgr.Wait.Wait()
}
