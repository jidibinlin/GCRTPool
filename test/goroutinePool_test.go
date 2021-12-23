package GCRTPool

import (
	"fmt"
	"testing"
	"time"

	GCRTPool "github.com/jidibinlin/GCRTPool/src"
)

// test ...
// func test(params ...interface{}) {
// 	//fmt.Println(params[0])
// 	for i := 0; i < 10000; i++ {
// 		GCRTPool.CoRun(testOne)
// 	}

// }

func testOne() {
	time.Sleep(time.Duration(10) * time.Millisecond)
}

func TestPerformance(t *testing.T) {
	mgr := GCRTPool.NewMgr()
	mgr.CreateCrts(2000, 100)
	fmt.Println("startTime", time.Now())
	for j := 0; j < threads; j++ {
		mgr.CoRun(func() {
			for i := 0; i < RunTimes; i++ {
				mgr.CoRun(testOne)
			}
		})
	}
	mgr.KillAllCoroutine(nil)
	mgr.Wait.Wait()
}
