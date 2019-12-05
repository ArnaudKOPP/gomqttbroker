package pool

const (
	// FNV-1a
	offset64 = uint64(14695981039346656037)
	prime64  = uint64(1099511628211)

	// Init64 is what 64 bits hash values should be initialized with.
	Init64 = offset64
)

//WorkerPool struct pool of workerq
type WorkerPool struct {
	maxWorkers  int
	taskQueue   []chan func()
	stoppedChan chan struct{}
}

// New new worker
func New(maxWorkers int) *WorkerPool {
	// There must be at least one worker.
	if maxWorkers < 1 {
		maxWorkers = 1
	}

	// taskQueue is unbuffered since items are always removed immediately.
	pool := &WorkerPool{
		taskQueue:   make([]chan func(), maxWorkers),
		maxWorkers:  maxWorkers,
		stoppedChan: make(chan struct{}),
	}
	// Start the task dispatcher.
	pool.dispatch()

	return pool
}

// HashString64 returns the hash of s.
func HashString64(s string) uint64 {
	return AddString64(Init64, s)
}

// AddString64 adds the hash of s to the precomputed hash value h.
func AddString64(h uint64, s string) uint64 {

	for len(s) >= 8 {
		h = (h ^ uint64(s[0])) * prime64
		h = (h ^ uint64(s[1])) * prime64
		h = (h ^ uint64(s[2])) * prime64
		h = (h ^ uint64(s[3])) * prime64
		h = (h ^ uint64(s[4])) * prime64
		h = (h ^ uint64(s[5])) * prime64
		h = (h ^ uint64(s[6])) * prime64
		h = (h ^ uint64(s[7])) * prime64
		s = s[8:]
	}

	if len(s) >= 4 {
		h = (h ^ uint64(s[0])) * prime64
		h = (h ^ uint64(s[1])) * prime64
		h = (h ^ uint64(s[2])) * prime64
		h = (h ^ uint64(s[3])) * prime64
		s = s[4:]
	}

	if len(s) >= 2 {
		h = (h ^ uint64(s[0])) * prime64
		h = (h ^ uint64(s[1])) * prime64
		s = s[2:]
	}

	if len(s) > 0 {
		h = (h ^ uint64(s[0])) * prime64
	}

	return h
}

// Submit new worker
func (p *WorkerPool) Submit(uid string, task func()) {
	idx := HashString64(uid) % uint64(p.maxWorkers)
	if task != nil {
		p.taskQueue[idx] <- task
	}
}

func (p *WorkerPool) dispatch() {
	for i := 0; i < p.maxWorkers; i++ {
		p.taskQueue[i] = make(chan func())
		go startWorker(p.taskQueue[i])
	}
}

func startWorker(taskChan chan func()) {
	go func() {
		var task func()
		var ok bool
		for {
			task, ok = <-taskChan
			if !ok {
				break
			}
			// Execute the task.
			task()
		}
	}()
}
