package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/vishalvivekm/cube/task"
	"github.com/vishalvivekm/cube/worker"
)

func main() {
	host := os.Getenv("CUBE_HOST")
	port, _ := strconv.Atoi(os.Getenv("CUBE_PORT"))
	fmt.Printf("starting cube worker at %s:%d\n", host, port)

	w := worker.Worker {
		Queue: *queue.New(),
		Db: make(map[uuid.UUID]*task.Task),
	}
	api := worker.Api{
		Address: host,
		Port: port,
		Worker: &w,
	}
	go runTasks(&w)
	go w.CollectStats()
	api.Start()
}

func runTasks(w *worker.Worker) {
	for {
		if w.Queue.Len() != 0 {
			res := w.RunTask()
			if res.Error != nil {
				log.Printf("error running task: %v\n", res.Error)
			}
		} else {
				log.Printf("no tasks to process currently. \n")
		}
		log.Println("sleeping for 15 seconds")
		time.Sleep(15 * time.Second)

	}
}