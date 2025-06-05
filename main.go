package main

import (
	"fmt"
	"workerpool/workerpool"
)

type PrintTask struct {
	num int
}

func (p PrintTask) Run() error {
	fmt.Printf("Task #%d running...\n", p.num)
	return nil
}

func main() {
	// usage example:

	var tasks []workerpool.Task
	for i := range 100 {
		tasks = append(tasks, PrintTask{i})
	}

	// create worker pool
	pool := workerpool.NewPool(len(tasks))

	// add some workers
	pool.AddWorker()
	pool.AddWorker()
	pool.AddWorker()
	pool.AddWorker()

	// submit tasks
	for i := range len(tasks) {
		pool.Submit(tasks[i])
	}

	// add & delete some workers

	// wait

	pool.Wait()

}
