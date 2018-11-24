package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	id int
}

var (
	wg sync.WaitGroup

	taskInChan  chan *Task
	taskOutChan chan *Task

	controlNumChan chan int
)

func init() {
	taskInChan = make(chan *Task, 10)
	taskOutChan = make(chan *Task, 10)

	controlNumChan = make(chan int, 10)
}

func worker(task *Task) {
	fmt.Println(task.id)
	task.id = task.id + 100
	taskOutChan <- task // return task
	time.Sleep(3 * time.Second)

	wg.Done()
	<-controlNumChan
}

func dispatch() {
	for {
		// select {
		// case task := <-taskInChan:
		// 	go worker(task)
		// }
		task := <-taskInChan
		go worker(task)
	}
}

func feedback() {
	for {
		// select {
		// case task := <-taskOutChan:
		// 	fmt.Println(task.id)
		// }
		task := <-taskOutChan
		fmt.Println(task.id)
	}
}

func main() {
	wg.Add(100) //task list len

	var ids []int
	i := 0
	for {
		ids = append(ids, i)
		i += 1
		if i >= 100 {
			break
		}
	}

	go dispatch()
	go feedback()
	for _, id := range ids {
		controlNumChan <- 1 // control gorutine num
		task := &Task{
			id: id,
		}
		taskInChan <- task // put task
	}

	wg.Wait() // wait all task finish
}
