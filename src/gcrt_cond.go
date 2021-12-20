package GCRTPool

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type GCRTCond struct {
	noCopy  noCopy
	L       sync.Locker
	notify  notifyList
	checker copyChecker
}

func newCond(l sync.Locker) *GCRTCond {
	return &GCRTCond{L: l}
}

// Wait ...
func (c *GCRTCond) wait(crt *crt) {
	c.checker.check()
	t := runtime_notifyListAdd(&c.notify)
	c.L.Unlock()
	GetMgr().readyCrtsMutex.Lock()
	GetMgr().readyCrts[crt.id] = crt
	GetMgr().readyCrtsMutex.Unlock()
	runtime_notifyListWait(&c.notify, t)
	c.L.Lock()
}

// Signal ...
func (c *GCRTCond) signal() {
	c.checker.check()
	runtime_notifyListNotifyOne(&c.notify)
}

func (c *GCRTCond) broadcast() {
	c.checker.check()
	runtime_notifyListNotifyAll(&c.notify)
}

type copyChecker uintptr

func (c *copyChecker) check() {
	if uintptr(*c) != uintptr(unsafe.Pointer(c)) &&
		!atomic.CompareAndSwapUintptr((*uintptr)(c), 0, uintptr(unsafe.Pointer(c))) &&
		uintptr(*c) != uintptr(unsafe.Pointer(c)) {
		panic("sync.Cond is copied")
	}
}

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
