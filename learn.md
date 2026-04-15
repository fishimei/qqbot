注意原子操作   bool -> atomic.Bool  来进行原子操作，避免竞态条件
锁操作        sync.Mutex  sync.RWMutex  sync.WaitGroup  写锁和读锁不能并存