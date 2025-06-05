package stringworkerpool

import (
	"context"
	"fmt"
	"sync"
)

type worker struct {
	id     int
	ctx    context.Context
	cancel context.CancelFunc
}

type StringWorkerPool struct {
	workers      map[int]*worker
	input        chan string
	lastWorkerId int

	wg     sync.WaitGroup
	mu     sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
}

// NewStringWorkerPool создает новый string worker pool с буферизированным каналом input.
func NewStringWorkerPool(capacity int) *StringWorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &StringWorkerPool{
		workers:      make(map[int]*worker),
		input:        make(chan string, capacity),
		lastWorkerId: 0,
		ctx:          ctx,
		cancel:       cancel,
	}
}

// AddWorker добавляет воркера и возвращает его id
func (pool *StringWorkerPool) AddWorker() int {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	pool.lastWorkerId++
	id := pool.lastWorkerId
	ctx, cancel := context.WithCancel(pool.ctx)
	w := &worker{id, ctx, cancel}
	pool.workers[id] = w

	pool.wg.Add(1)
	go func(w *worker) {
		defer pool.wg.Done()
		fmt.Printf("StringWorker #%d started\n", w.id)
		for {
			select {
			case msg, ok := <-pool.input:
				if !ok {
					fmt.Printf("StringWorker #%d stopped\n", w.id)
					return
				}
				fmt.Printf("[Worker %d] %s\n", w.id, msg)
			case <-w.ctx.Done():
				fmt.Printf("StringWorker #%d stopped\n", w.id)
				return
			}
		}
	}(w)
	return id
}

// RemoveWorker останавливает и удаляет воркера по id
func (pool *StringWorkerPool) RemoveWorker(id int) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	if w, ok := pool.workers[id]; ok {
		w.cancel()
		delete(pool.workers, id)
	}
}

// Stop завершает работу всех воркеров и закрывает канал
func (pool *StringWorkerPool) Stop() {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	pool.cancel()
	for _, w := range pool.workers {
		w.cancel()
	}
	pool.workers = make(map[int]*worker)
	close(pool.input)
}

// Input возвращает канал для ввода строки
func (pool *StringWorkerPool) Input() chan<- string {
	return pool.input
}

// Wait ожидает завершения всех воркеров
func (pool *StringWorkerPool) Wait() {
	pool.wg.Wait()
}
