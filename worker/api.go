package worker

import (
	"fmt"
	"net/http"
)

type Api struct {
	Address string
	Port int
	Worker *Worker
	Router *http.ServeMux
}
type ErrResponse struct {
	HTTPStatuscode int
	Message string
}
func (a *Api) initRouter() {
	r := http.NewServeMux()
	r.HandleFunc("GET /tasks", a.GetTaskHandler)
	r.HandleFunc("POST /tasks", a.StartTaskHandler)
	r.HandleFunc("DELETE /tasks/{taskID}", a.StopTaskHandler)
	a.Router = r
}
func (a *Api) Start() {
	a.initRouter()
	http.ListenAndServe(fmt.Sprintf("%s:%d", a.Address, a.Port), a.Router)
}