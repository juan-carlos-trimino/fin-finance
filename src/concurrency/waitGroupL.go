package concurrency

/***
***/
type WaitGrpL struct {
  sp *Semaphore
}

func NewWaitGrpL(count int) *WaitGrpL {
  return &WaitGrpL{
    sp: NewSemaphore(1 - count),
  }
}

func (wg *WaitGrpL) Wait() {
  wg.sp.Acquire()
}

func (wg *WaitGrpL) Done() {
  wg.sp.Release()
}
