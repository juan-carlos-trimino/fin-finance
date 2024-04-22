package concurrency

import "sync"

/***
Whereas waitgroups are great for synchronizing after a task has been completed, barriers provide
the ability to synchronize groups of goroutines at specific points in the code. Barriers are
different from waitgroups in that they combine the waitgroups' Done() and Wait() functions together
into one atomic call.
Notes:
(1) Barriers suspend a goroutine when the goroutine calls Wait() until all of the goroutines
    participating in the barrier also call Wait().
(2) When all of the goroutines participating in the barrier call Wait(), all of the suspended
    goroutines on the barrier are resumed.
(3) Barriers can be reused multiple times.
***/
type Barrier struct {
  //The counter for the barrier.
  count int
  //Counter representing the number of currently suspended goroutines.
  waiting int
  cond *sync.Cond
}

func NewBarrier(count int) *Barrier {
  return &Barrier{
    count: count,
    waiting: 0,
    cond: sync.NewCond(&sync.Mutex{}),
  }
}

func (b *Barrier) Wait() {
  //Protect the waiting variable.
  b.cond.L.Lock()
  b.waiting++
  //If waiting has reached the barrier count, reset waiting and broadcast on the condition
  //variable.
  if b.waiting == b.count {
    //Depending on the implementation, barriers can be reused multiple times (cyclic barriers).
    b.waiting = 0
    b.cond.Broadcast()
  } else {
    //If waiting hasn't reached the barrier count, wait on the condition variable.
    b.cond.Wait()
  }
  b.cond.L.Unlock()
}
