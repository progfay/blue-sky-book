package queue

import (
	"context"
	"sync"
)

type Queue struct {
	queue         []string
	mu            sync.Mutex
	idlingWorkers chan worker
	isRunning     bool
	wg            sync.WaitGroup
	ctx           context.Context
	cancel        func()
}

type worker struct {
	ID      int64
	Handler func(string)
}

func (w *worker) Run(job string) {
	w.Handler(job)
}

func NewQueue(ctx context.Context, maxJobCount int, handler func(string)) Queue {
	idlingWorkers := make(chan worker, maxJobCount)
	for i := 0; i < maxJobCount; i++ {
		idlingWorkers <- worker{
			ID:      int64(i),
			Handler: handler,
		}
	}

	childCtx, cancel := context.WithCancel(ctx)

	return Queue{
		queue:         make([]string, 0),
		mu:            sync.Mutex{},
		idlingWorkers: idlingWorkers,
		isRunning:     false,
		wg:            sync.WaitGroup{},
		ctx:           childCtx,
		cancel:        cancel,
	}
}

func (q *Queue) Start() {
	if q.isRunning {
		return
	}

	q.isRunning = true

	for {
		select {
		case <-q.ctx.Done():
			q.wg.Wait()
			return
		case w := <-q.idlingWorkers:
			if len(q.queue) == 0 {
				continue
			}
			func() {
				q.mu.Lock()
				defer q.mu.Unlock()

				job := q.queue[0]
				q.queue = q.queue[1:]

				if len(q.queue) == 0 {
					q.wg.Wait()
					q.isRunning = false
					q.cancel()
				}

				q.wg.Add(1)
				go func(job string) {
					w.Run(job)
					q.wg.Done()
					q.idlingWorkers <- w
				}(job)
			}()
		}
	}
}

func (q *Queue) Add(job string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue = append(q.queue, job)
}
