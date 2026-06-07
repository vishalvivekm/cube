package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/vishalvivekm/cube/task"
	"github.com/vishalvivekm/cube/worker"
	"github.com/vishalvivekm/cube/manager"
)

func main() {
	whost := os.Getenv("CUBE_WORKER_HOST")
	wport, _ := strconv.Atoi(os.Getenv("CUBE_WORKER_PORT"))
	mhost := os.Getenv("CUBE_MANAGER_HOST")
	mport, _ := strconv.Atoi(os.Getenv("CUBE_MANAGER_PORT"))

	fmt.Printf("starting cube worker at %s:%d\n", whost, wport)

	w := worker.Worker {
		Queue: *queue.New(),
		Db: make(map[uuid.UUID]*task.Task),
	}
	wapi := worker.Api{
		Address: whost,
		Port: wport,
		Worker: &w,
	}
	go w.RunTasks()
	go w.CollectStats()
	go wapi.Start()

	fmt.Printf("starting cube manager at %s:%d\n", mhost, mport)
	
	workers := []string{fmt.Sprintf("%s:%d", whost, wport)}
	m := manager.New(workers)
	mapi := manager.Api{Address: mhost, Port: mport, Manager: m}
 
	go m.ProcessTasks()
	go m.UpdateTasks()
	mapi.Start()
}