package manager

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/google/uuid"
	"github.com/vishalvivekm/cube/task"
	"time"

)

func (a *Api) StartTaskHandler(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	te := task.TaskEvent{}
	if err := d.Decode(&te); err != nil {
		msg := fmt.Sprintf("error unmarshaling req body; %v\n", err)
		log.Printf("msg")
		w.WriteHeader(400)
		e := ErrResponse{
			HTTPStatuscode: 400,
			Message: msg,	
		}
		json.NewEncoder(w).Encode(e)
		return
	}
	a.Manager.AddTask(te)	
	log.Printf("added task %v\n", te.Task.ID)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(te.Task)
}

func(m *Manager) GetTasks() []*task.Task {
	tasks := []*task.Task{}
	for _, t := range m.TaskDb{
		tasks = append(tasks, t)
	}
	return tasks
}


func (a *Api) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(a.Manager.GetTasks())
}

func (a *Api) StopTaskHandler(w http.ResponseWriter, r *http.Request) {
	d := r.PathValue("taskID")
	if d == "" {
		log.Printf("no task id passed in req\n")
		w.WriteHeader(http.StatusBadRequest)
	}
	tID, _ := uuid.Parse(d)
	taskStop, ok := a.Manager.TaskDb[tID]	
	if !ok {
		log.Printf("no task with ID %v found\n",tID)
		w.WriteHeader(http.StatusNotFound)
	}
	te := task.TaskEvent{
		ID: uuid.New(), 
		State: task.Completed, 
		Timestamp: time.Now(),
	}
	taskCopy := *taskStop // make a copy of taskToStop, and change state of the copy
	taskCopy.State = task.Completed
	te.Task = taskCopy
	a.Manager.AddTask(te)
	log.Printf("added task event %v to stop task %v\n", te.ID, taskStop.ID)
	w.WriteHeader(http.StatusNoContent)
}

