package concurrency

import (
  "sync"
)

type WaitGrp struct {
  count int
  cond *sync.Cond
}

func NewWaitGrp() *WaitGrp {
  return &WaitGrp{
    count: 0,
    cond: sync.NewCond(&sync.Mutex{}),
  }
}

func (wg *WaitGrp) Add(delta int) {
  wg.cond.L.Lock()
  wg.count += delta
  wg.cond.L.Unlock()
}

func (wg *WaitGrp) Wait() {
  wg.cond.L.Lock()
  for wg.count > 0 {
    wg.cond.Wait()
  }
  wg.cond.L.Unlock()
}

func (wg *WaitGrp) Done() {
  wg.cond.L.Lock()
  wg.count--
  if wg.count == 0 {
    wg.cond.Broadcast()
  }
  wg.cond.L.Unlock()
}