package concurrency

import (
  "sync"
)

/***
Semaphores have the ability to limit concurrency on a shared resource to a fixed number of
concurrent goroutines. Furthermore, when used to signal, semaphores have the added advantage that
the signal is stored if the goroutine is not yet waiting for it.
***/
type Semaphore struct {
  //Counter for the semaphore.
  count int
  //Condition variable used for waiting when the counter is zero (0).
  cond *sync.Cond
}

func NewSemaphore(c int) *Semaphore {
  return &Semaphore{
    //Initial count.
    count: c,
    cond: sync.NewCond(&sync.Mutex{}),
  }
}

func (rw *Semaphore) Acquire() {
  rw.cond.L.Lock()
  for rw.count < 1 {
    rw.cond.Wait()
  }
  rw.count--
  rw.cond.L.Unlock()
}

func (rw *Semaphore) Release() {
  rw.cond.L.Lock()
  rw.count++
  //Use Signal() instead of Broadcast() since rw.count was incremented by one; hence, only one
  //goroutine should be released.
  rw.cond.Signal()
  rw.cond.L.Unlock()
}
