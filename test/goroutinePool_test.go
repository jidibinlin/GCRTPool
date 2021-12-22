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
	for i := 0; i < 1000; i++ {

	}

	//time.Sleep(1 * time.Millisecond)

}

func TestPerformance(t *testing.T) {
	mgr := GCRTPool.GetMgr()
	mgr.CreateCrts(2, 2)
	//time.Sleep(1 * time.Second)
	beforeTest := time.Now()
	// for i := 0; i < 100; i++ {
	// 	//testCoroutine100000()
	// 	//fmt.Println("run", i)
	// 	GCRTPool.CoRun(test, 1)
	// }
	i := 1
	GCRTPool.CoRun(func(params []interface{}) {
		fmt.Println(i)
		i += 1
		fmt.Println(i)
	}, nil)

	mgr.KillAllCoroutine(nil)
	mgr.Wait.Wait()

	afterTest := time.Now()

	fmt.Println(afterTest, beforeTest)
}
