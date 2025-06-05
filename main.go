package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"workerpool/stringworkerpool"
	"workerpool/workerpool"
)

type PrintTask struct {
	num int
}

func (p PrintTask) Run() error {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // random time from 0 to 1 sec
	fmt.Printf("Task #%d is running...\n", p.num)
	return nil
}

func main() {

	// Пример использования Worker Pool с тасками

	var tasks []workerpool.Task
	for i := range 100 {
		tasks = append(tasks, PrintTask{i})
	}

	pool := workerpool.NewPool(len(tasks))

	pool.AddWorker()
	second := pool.AddWorker()
	pool.AddWorker()
	pool.AddWorker()

	for i := range len(tasks) {
		pool.Submit(tasks[i])
	}

	pool.RemoveWorker(second)

	pool.Wait()

	// stop all workers
	// time.Sleep(5 * time.Second)
	// pool.Stop()

	// Пример использования StringWorkerPool
	var strings []string
	for i := range 20 {
		strings = append(strings, "string #"+strconv.Itoa(i))
	}

	sPool := stringworkerpool.NewStringWorkerPool(len(strings))
	sPool.AddWorker()
	second = sPool.AddWorker()
	sPool.AddWorker()
	sPool.AddWorker()

	for _, s := range strings {
		sPool.Input() <- s
	}
	close(sPool.Input())

	sPool.RemoveWorker(second)

	sPool.Wait()

}
