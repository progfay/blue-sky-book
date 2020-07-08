package queue

import "sync"

type Queue struct {
	queue         []func()
	mu            sync.Mutex
	idlingWorkers chan worker
	isRunning     bool
}

type worker struct{}

func (w *worker) Run(job func()) {
	job()
}

func NewQueue(maxJobCount int) Queue {
	idlingWorkers := make(chan worker, maxJobCount)
	for i := 0; i < maxJobCount; i++ {
		idlingWorkers <- worker{}
	}

	return Queue{
		queue:         make([]func(), 0),
		mu:            sync.Mutex{},
		idlingWorkers: idlingWorkers,
		isRunning:     false,
	}
}

func (q *Queue) Start() {
	if q.isRunning {
		return
	}

	q.isRunning = true

	for {
		select {
		case w := <-q.idlingWorkers:
			q.mu.Lock()
			defer q.mu.Unlock()
			if len(q.queue) == 0 {
				q.isRunning = false
				return
			}

			job := q.queue[0]
			q.queue = q.queue[1:]

			go func(job func()) {
				w.Run(job)
				q.idlingWorkers <- w
			}(job)
		}
	}
}

func (q *Queue) Add(job func()) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.queue = append(q.queue, job)
}
