package main

// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"strconv"
// 	"time"

// 	"github.com/golang-collections/collections/queue"
// 	"github.com/google/uuid"
// 	"github.com/vishalvivekm/cube/task"
// 	"github.com/vishalvivekm/cube/worker"
// 	"github.com/vishalvivekm/cube/manager"
// )

// func main() {
// 	host := os.Getenv("CUBE_HOST")
// 	port, _ := strconv.Atoi(os.Getenv("CUBE_PORT"))
// 	fmt.Printf("starting cube worker at %s:%d\n", host, port)

// 	w := worker.Worker {
// 		Queue: *queue.New(),
// 		Db: make(map[uuid.UUID]*task.Task),
// 	}
// 	api := worker.Api{
// 		Address: host,
// 		Port: port,
// 		Worker: &w,
// 	}
// 	go runTasks(&w)
// 	go w.CollectStats()
// 	go api.Start()
// 	workers := []string{fmt.Sprintf("%s:%d", host, port)}
// 	m := manager.New(workers)
// 	for i := 0; i < 3; i++ {
// 		t := task.Task {
// 			ID: uuid.New(), 
// 			Name: fmt.Sprintf("test-container-%d", i),
// 			State: task.Scheduled, 
// 			Image: "vishalvivekm/to-do-app:latest", 
// 		}
// 		te := task.TaskEvent {
// 			ID: uuid.New(),
// 			State: task.Running, 
// 			Task: t,
// 		}
// 		m.AddTask(te)
// 		m.SendWork()

// 		go func(){
// 			for {
// 				fmt.Printf("[Manager] Updating tasks from %d workers\n", len(m.Workers))
// 				m.UpdateTasks()
// 				time.Sleep(15 * time.Second)
// 			}
// 		}()
// 		for _, t := range m.TaskDb{
// 			fmt.Printf("[Manager] Task: id: %s, state: %d\n", t.ID, t.State)
// 			time.Sleep((15 * time.Second))
// 		}
// 	}
// }

// func runTasks(w *worker.Worker) {
// 	for {
// 		if w.Queue.Len() != 0 {
// 			res := w.RunTask() // will break here
// 			if res.Error != nil {
// 				log.Printf("error running task: %v\n", res.Error)
// 			}
// 		} else {
// 				log.Printf("no tasks to process currently. \n")
// 		}
// 		log.Println("sleeping for 15 seconds")
// 		time.Sleep(15 * time.Second)

// 	}
// }