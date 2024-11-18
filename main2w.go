package main

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/vishalvivekm/cube/task"
	"github.com/vishalvivekm/cube/worker"
)

func main() {
	db := make(map[uuid.UUID]*task.Task)
	w := worker.Worker{
		Queue: *queue.New(),
		Db: db,
	}

	t := task.Task{
		ID: uuid.New(),
		Name: "test-container-1",
		State: task.Scheduled,
		// Image: "strm/helloworld-http",
		Image: "vishalvivekm/to-do-app",
	}

	fmt.Println("starting task")
	w.AddTask(t)
	res := w.RunTask()
	if res.Error != nil {
		log.Fatalln(res.Error)
	}

	t.ContainerID = res.ContainerId
	fmt.Printf("task %s in container %s\n", t.ID, t.ContainerID)
	fmt.Println("sleeping...")
	time.Sleep(time.Second * 60)

	fmt.Printf("stopping task %s\n", t.ID)
	t.State = task.Completed
	w.AddTask(t)
	res = w.RunTask()
	if res.Error != nil {
		log.Fatalln(res.Error)
	}
}