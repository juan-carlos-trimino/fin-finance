package concurrency

import (
  "sync"
)

/***
This implementation does not have the limitations described in the WaitGrpL; i.e., the count of the
waitgroup can be changed after creation and more than one goroutine can be suspended and unblocked
on the Wait() function.
***/
type WaitGrp struct {
  //The waitgroup count; Go initializes it to 0 by default.
  count int
  cond *sync.Cond
}

func NewWaitGrp() *WaitGrp {
  return &WaitGrp{
    count: 0,
    //Initialize the condition variable with a mutex.
    cond: sync.NewCond(&sync.Mutex{}),
  }
}

func (wg *WaitGrp) Add(delta int) {
  wg.cond.L.Lock()  //Protect count.
  wg.count += delta
  wg.cond.L.Unlock()
}

func (wg *WaitGrp) Wait() {
  //Protect the read of count with a mutex lock on the condition variable.
  wg.cond.L.Lock()
  for wg.count > 0 {
    //Wait and release the mutex atomically while count is greater than 0.
    wg.cond.Wait()
  }
  wg.cond.L.Unlock()
}

func (wg *WaitGrp) Done() {
  wg.cond.L.Lock()  //Protect count.
  wg.count--
  //If it's the last goroutine to be done in the waitgroup, it broadcasts on the condition
  //variable.
  if wg.count == 0 {
    wg.cond.Broadcast()
  }
  wg.cond.L.Unlock()
}
