package worker

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-collections/collections/queue"
	"github.com/google/uuid"
	"github.com/vishalvivekm/cube/task"
)

type Worker struct {
	Name string
	Queue queue.Queue
	Db map[uuid.UUID]*task.Task
	TaskCount int // number of tasks worker's been assigned
	Stats *Stats
}


	func (w *Worker) CollectStats() {
		for {
			log.Println("collecting stats")
			w.Stats = GetStats()
			w.Stats.TaskCount = w.TaskCount
			time.Sleep(20 * time.Second)
		}
	}
	func (w *Worker) RunTask() task.DockerResult {
	t := w.Queue.Dequeue()
	if t == nil {
		log.Println("no tasks in the queue")
		return task.DockerResult{Error: nil}
	}
	taskQueued := t.(task.Task)
	fmt.Printf("found task in queue: %v:\n", taskQueued)
	taskPersisted := w.Db[taskQueued.ID]
	if taskPersisted == nil { // first time seeing this task :)
		// persist it in the db
		taskPersisted = &taskQueued
		w.Db[taskPersisted.ID] = &taskQueued
	}

	var res task.DockerResult
	if task.ValidStateTransition(
		taskPersisted.State, taskQueued.State) {
			switch taskQueued.State {
			case task.Scheduled:
				res = w.StartTask(taskQueued)
			case task.Completed:
				res = w.StopTask(taskQueued)
			default:
				res.Error = errors.New("we should not get here")
			}
		} else {
			err := fmt.Errorf("invalid transition from %v to %v", 
					taskPersisted.State, taskQueued.State,
			) 
			res.Error = err
		}
	return res
	}

	func (w *Worker) StartTask(t task.Task) task.DockerResult {

		t.StartTime = time.Now().UTC()
		config := task.NewConfig(&t)
		d := task.NewDocker(config)
		res := d.Run()
		if res.Error != nil {
			log.Printf("error running task %v: %v\n", t.ID, res.Error)
			t.State = task.Failed
			w.Db[t.ID] = &t
			return res
		}
		t.ContainerID = res.ContainerId
		t.State = task.Running
		w.Db[t.ID] = &t

		return res
		
	}
	
	func (w *Worker) StopTask(t task.Task) task.DockerResult {
		config := task.NewConfig(&t)
		d := task.NewDocker(config)

		res := d.Stop(t.ContainerID)
		if res.Error != nil {
			log.Printf("error stopping container %v: %v\n", t.ContainerID,
			res.Error)
		}
		t.FinishTime = time.Now().UTC()
		t.State = task.Completed
		w.Db[t.ID] = &t
		log.Printf("Stopped and removed container %v for task %v\n",t.ContainerID,
	t.ID)
	
	return res

	}
	// Add task to the queue
	func (w *Worker) AddTask(t task.Task) {
		w.Queue.Enqueue(t)
	}

	// return all the tasks in worker's db
	func (w *Worker) GetTasks() []*task.Task {
	tasks := make([]*task.Task, 0)
	for _, task := range w.Db {
		tasks = append(tasks, task)
	}	
	return tasks
	}
