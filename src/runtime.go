// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package GCRTPool

import "unsafe"

// defined in package runtime

// Semacquire waits until *s > 0 and then atomically decrements it.
// It is intended as a simple sleep primitive for use by the synchronization
// library and should not be used directly.
//go:linkname runtime_Semacquire sync.runtime_Semacquire
func runtime_Semacquire(s *uint32)

// SemacquireMutex is like Semacquire, but for profiling contended Mutexes.
// If lifo is true, queue waiter at the head of wait queue.
// skipframes is the number of frames to omit during tracing, counting from
// runtime_SemacquireMutex's caller.
//go:linkname runtime_SemacquireMutex sync.runtime_SemacquireMutex
func runtime_SemacquireMutex(s *uint32, lifo bool, skipframes int)

// Semrelease atomically increments *s and notifies a waiting goroutine
// if one is blocked in Semacquire.
// It is intended as a simple wakeup primitive for use by the synchronization
// library and should not be used directly.
// If handoff is true, pass count directly to the first waiter.
// skipframes is the number of frames to omit during tracing, counting from
// runtime_Semrelease's caller.
//go:linkname runtime_Semrelease sync.runtime_Semrelease
func runtime_Semrelease(s *uint32, handoff bool, skipframes int)

// See runtime/sema.go for documentation.
//go:linkname runtime_notifyListAdd sync.runtime_notifyListAdd
func runtime_notifyListAdd(l *notifyList) uint32

// See runtime/sema.go for documentation.
//go:linkname runtime_notifyListWait sync.runtime_notifyListWait
func runtime_notifyListWait(l *notifyList, t uint32)

// See runtime/sema.go for documentation.
//go:linkname runtime_notifyListNotifyAll sync.runtime_notifyListNotifyAll
func runtime_notifyListNotifyAll(l *notifyList)

// See runtime/sema.go for documentation.
//go:linkname runtime_notifyListNotifyOne sync.runtime_notifyListNotifyOne
func runtime_notifyListNotifyOne(l *notifyList)

// Ensure that sync and runtime agree on size of notifyList.
//go:linkname runtime_notifyListCheck sync.runtime_notifyListCheck
func runtime_notifyListCheck(size uintptr)
func init() {
	var n notifyList
	runtime_notifyListCheck(unsafe.Sizeof(n))
}

// Active spinning runtime support.
// runtime_canSpin reports whether spinning makes sense at the moment.
//go:linkname runtime_canSpin sync.runtime_canSpin
func runtime_canSpin(i int) bool

// runtime_doSpin does active spinning.
//go:linkname runtime_doSpin
func runtime_doSpin()

//go:linkname runtime_nanotime sync.runtime_nanotime
func runtime_nanotime() int64
