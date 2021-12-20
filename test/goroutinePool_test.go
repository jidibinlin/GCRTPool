package GCRTPool

import (
	"fmt"
	"testing"
	"time"

	GCRTPool "github.com/jidibinlin/GCRTPool/src"
)

// test ...
func test(params []interface{}) {
	//fmt.Println(params[0])
	for i := 0; i < 100000; i++ {
		GCRTPool.CoRun(testCoroutine100000, 1)
	}

}

// testCoroutine100000 ...
func testCoroutine100000(params []interface{}) {
}

func TestPerformance(t *testing.T) {
	mgr := GCRTPool.GetMgr()
	mgr.CreateCrts(100000)
	//time.Sleep(1 * time.Second)
	beforeTest := time.Now()
	for i := 0; i < 100; i++ {
		//testCoroutine100000()
		//fmt.Println("run", i)
		GCRTPool.CoRun(test, 1)
	}

	GCRTPool.CoRun(testCoroutine100000)

	mgr.KillAllCoroutine(nil)

	mgr.Wait.Wait()

	afterTest := time.Now()

	fmt.Println(afterTest, beforeTest)
}
