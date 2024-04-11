package concurrency

/***
Implementation of a waitgroup using a semaphore.
Limitations:
(1) The size of the waitgroup is specified at creation, and it cannot be changed.
(2) Only one goroutine can wait on the waitgroup. If multiple goroutines call the Wait() function,
    only one will resume because the count of the semaphore is incremented by 1.
***/
type WaitGrpL struct {
  sp *Semaphore
}

func NewWaitGrpL(count int) *WaitGrpL {
  /***
  When creating the waitgroup, initialize the semaphore to a count of (1 - n), where n is the size
  of the waitgroup. This will force the Wait() function to block until the count is increased n
  times, from (1 - n) to 1.
  ***/
  return &WaitGrpL{
    sp: NewSemaphore(1 - count),
  }
}

func (wg *WaitGrpL) Wait() {
  //Calling Acquire() on the semaphore.
  wg.sp.Acquire()
}

func (wg *WaitGrpL) Done() {
  //Calling Release() on the semaphore.
  wg.sp.Release()
}
