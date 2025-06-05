package workerpool

import (
	"context"
	"fmt"
	"sync"
)

type Task interface {
	Run() error
}

// worker - структура одного воркера
type worker struct {
	id     int
	ctx    context.Context
	cancel context.CancelFunc
}

// Pool - основная структура ..
type Pool struct {
	workers      map[int]*worker // map со всеми воркерами
	tasks        chan Task       // канал с задачами для выполнения
	lastWorkerId int             // id последнего добавленного воркера

	wg     sync.WaitGroup
	mu     sync.Mutex         // мьютекс для безопасной работы с общей мапой workers
	ctx    context.Context    // общий контекст
	cancel context.CancelFunc // общая отмена
}

// NewPool - создаёт новый воркер пул
func NewPool(capacity int) *Pool {

	ctx, cancel := context.WithCancel(context.Background())

	return &Pool{
		workers:      make(map[int]*worker),
		tasks:        make(chan Task, capacity),
		lastWorkerId: 0,
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Submit - отправляет задачу в пул для обработки
func (pool *Pool) Submit(task Task) {
	pool.wg.Add(1)
	pool.tasks <- task
}

// AddWorker - добавляет одного воркера и запускает его в отдельной горутине
// возвращает id созданного воркера
func (pool *Pool) AddWorker() int {
	// чтобы не было гонки
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// создадим нового воркера
	pool.lastWorkerId++
	id := pool.lastWorkerId
	ctx, cancel := context.WithCancel(pool.ctx)
	w := &worker{id, ctx, cancel}
	// добавим в мапу
	pool.workers[id] = w

	// запустим
	go func(w *worker) {
		fmt.Printf("Worker #%d started\n", w.id)
		for {
			select {
			case task := <-pool.tasks:
				if err := task.Run(); err != nil {
					fmt.Printf("Worker #%d: error %v occured while executing task\n", w.id, err)
				}
				pool.wg.Done()
			case <-w.ctx.Done():
				fmt.Printf("Worker #%d finished	its work\n", w.id)
				return
			}
		}
	}(w)

	return id
}

func (pool *Pool) RemoveWorker(id int) {
	// чтобы не было гонки
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// Если есть такой воркер, то остановим его работу
	// и удалим из мапы
	if worker, ok := pool.workers[id]; ok {
		worker.cancel()
		delete(pool.workers, id)
	}
}

// Stop - остановка и удаление всех воркеров
func (pool *Pool) Stop() {
	// чтобы не было гонки
	pool.mu.Lock()
	defer pool.mu.Unlock()

	// общая отмена и пересоздать мапу
	pool.cancel()
	pool.workers = make(map[int]*worker)
}

func (pool *Pool) Wait() {
	pool.wg.Wait()
}
