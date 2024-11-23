package worker
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/google/uuid"
	"github.com/vishalvivekm/cube/task"

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
	a.Worker.AddTask(te.Task)	
	log.Printf("added task %v\n", te.Task.ID)
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(te.Task)
}



func (a *Api) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(a.Worker.GetTasks())
}

func (a *Api) StopTaskHandler(w http.ResponseWriter, r *http.Request) {
	d := r.PathValue("taskID")
	if d == "" {
		log.Printf("no task id passed in req\n")
		w.WriteHeader(http.StatusBadRequest)
	}
	tID, _ := uuid.Parse(d)
	taskStop, ok := a.Worker.Db[tID]	
	if !ok {
		log.Printf("no task with ID %v found\n",tID)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	taskCopy := *taskStop // make a copy of taskToStop, and change state of the copy
	taskCopy.State = task.Completed
	a.Worker.AddTask(taskCopy)
	log.Printf("added task %v to stop container %v\n", taskStop.ID, taskStop.ContainerID)
	w.WriteHeader(http.StatusNoContent)
}