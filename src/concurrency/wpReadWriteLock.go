package concurrency

import (
  "sync"
)

/***
Write-preferred readers-writer lock.
***/
type WPReadWriteLock struct {
  //Indicate if a writer goroutine is holding the write lock.
  writerActive bool
  //The number of reader goroutines currently holding the read lock.
  readersCounter int
  //The number of writer goroutines currently waiting.
  writersWaiting int
  /***
  Condition variables work together with mutexes and provide the ability to suspend the current
  goroutine until there is a signal that a particular condition has changed.

  To create a new Go condition variable requires a Locker (interface), which defines two functions:
  Lock() and Unlock(). To use Go's condition variable, a type that implements these two functions
  is required, and a mutex is one such type.
  ***/
  cond *sync.Cond
}

/***
Initialize a new WPReadWriteLock with a new condition variable and associated mutex.
***/
func NewWPReadWriteLock() *WPReadWriteLock {
  return &WPReadWriteLock{
    writerActive: false,
    readersCounter: 0,
    writersWaiting: 0,
    cond: sync.NewCond(&sync.Mutex{}),
  }
}

func (rw *WPReadWriteLock) ReadLock() {
  rw.cond.L.Lock()  //Acquire mutex.
  for rw.writersWaiting > 0 || rw.writerActive {
    /***
    (1) When there are multiple goroutines suspended on a condition variable's Wait(), Signal()
        will arbitrarily wake up one of these goroutines. On the other hand, Broadcast() will wake
        up all goroutines that are suspended on a Wait().
    (2) Whenever a waiting goroutine receives a signal or broadcast, it will try to reacquire the
        mutex. If another goroutine is holding the mutex, the goroutine will remain suspended until
        the mutex becomes available.
    (3) The Wait() function releases the mutex and suspends the goroutine in an atomic manner. This
        means that another goroutine cannot come in between these two operations; i.e., acquire the
        lock and call the Signal() function before the goroutine calling Wait() has been suspended.
    (4) If Signal() or Broadcast() is called and no goroutines are suspended on a Wait(), the
        signal or broadcast is missed.
    ***/
    rw.cond.Wait()
  }
  rw.readersCounter++
  rw.cond.L.Unlock()  //Release mutex.
}

func (rw *WPReadWriteLock) ReadUnlock() {
  rw.cond.L.Lock()
  rw.readersCounter--
  //If last reader, send a broadcast.
  if rw.readersCounter == 0 {
    rw.cond.Broadcast()
  }
  rw.cond.L.Unlock()
}

func (rw *WPReadWriteLock) WriteLock() {
  rw.cond.L.Lock()
  rw.writersWaiting++
  for rw.readersCounter > 0 || rw.writerActive {
    rw.cond.Wait()
  }
  rw.writersWaiting--
  rw.writerActive = true
  rw.cond.L.Unlock()
}

func (rw *WPReadWriteLock) WriteUnlock() {
  rw.cond.L.Lock()
  rw.writerActive = false
  rw.cond.Broadcast()
  rw.cond.L.Unlock()
}
