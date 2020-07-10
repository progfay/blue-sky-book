package queue

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type Queue struct {
	queue         []string
	mu            sync.Mutex
	idlingWorkers chan worker
	isRunning     bool
	wg            sync.WaitGroup
	ctx           context.Context
	cancel        func()
	doneJobCount  int64
	allJobCount   int64
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
		doneJobCount:  0,
		allJobCount:   0,
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
			q.isRunning = false
			return
		case w := <-q.idlingWorkers:
			func() {
				q.mu.Lock()
				defer q.mu.Unlock()

				if len(q.queue) == 0 {
					q.cancel()
					return
				}

				job := q.queue[0]
				q.queue = q.queue[1:]

				q.wg.Add(1)
				go func(job string) {
					w.Run(job)
					atomic.AddInt64(&q.doneJobCount, 1)
					q.idlingWorkers <- w
					width, _, _ := terminal.GetSize(syscall.Stdin)
					doneJobCount := int(atomic.LoadInt64(&q.doneJobCount))
					allJobCount := int(atomic.LoadInt64(&q.allJobCount))
					done := 0
					if allJobCount > 0 {
						done = (width - 16) * doneJobCount / allJobCount
					}
					fmt.Fprintf(os.Stderr, "\r[\x1b[7m%s\x1b[0m%s] %5d / %5d", strings.Repeat(" ", done), strings.Repeat(" ", width-done-16), doneJobCount, allJobCount)
					q.wg.Done()
				}(job)
			}()
		}
	}
}

func (q *Queue) Add(job string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue = append(q.queue, job)
	atomic.AddInt64(&q.allJobCount, 1)
}
