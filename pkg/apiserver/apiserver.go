package apiserver

import (
	"encoding/json"
	"fmt"
	"github.com/Alex27Khalupka/Go-course-task/pkg/model"
	"github.com/Alex27Khalupka/Go-course-task/pkg/service"
	"github.com/Alex27Khalupka/Go-course-task/pkg/store"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type APIServer struct {
	config *Config
	router *mux.Router
	Store  *store.Store
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		router: mux.NewRouter(),
	}
}

// func Start() configures store and router
func (s *APIServer) Start() error {

	fmt.Println("starting api server")

	s.configureRouter()

	if err := s.configureStore(); err != nil {
		return err
	}

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

// func configureRouter sets handlers for specific urls
func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/groups", s.handleGetGroups).Methods(http.MethodGet)
	s.router.HandleFunc("/tasks", s.handleGetTasks).Methods(http.MethodGet)
	s.router.HandleFunc("/groups", s.handlePostGroups).Methods(http.MethodPost)
	s.router.HandleFunc("/tasks", s.handlePostTasks).Methods(http.MethodPost)
	s.router.HandleFunc("/timeframes", s.handlePostTimeFrames).Methods(http.MethodPost)
	s.router.HandleFunc("/groups/{id}", s.handlePutGroups).Methods(http.MethodPut)
	s.router.HandleFunc("/tasks/{id}", s.handlePutTasks).Methods(http.MethodPut)
	s.router.HandleFunc("/groups/{id}", s.handleDeleteGroups).Methods(http.MethodDelete)
	s.router.HandleFunc("/tasks/{id}", s.handleDeleteTasks).Methods(http.MethodDelete)
	s.router.HandleFunc("/timeframes/{id}", s.handleDeleteTimeFrames).Methods(http.MethodDelete)
}

func (s *APIServer) configureStore() error {
	st := store.New(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}
	s.Store = st

	return nil
}

// func handleGetGroups handles GET request for /groups
func (s *APIServer) handleGetGroups(w http.ResponseWriter, req *http.Request) {
	if err := s.Store.Open(); err != nil {
		log.Fatal(err)
		return
	}

	groups := model.ResponseGroups{Groups: service.GetGroups(s.Store.GetDB())}

	jsonResponse, err := json.Marshal(groups)
	if err != nil {
		log.Fatal(err)
		return
	}

	if _, err = w.Write(jsonResponse); err != nil {
		log.Fatal(err)
		return
	}
}

// func handleGetGroups handles GET request for /tasks
func (s *APIServer) handleGetTasks(w http.ResponseWriter, req *http.Request) {
	if err := s.Store.Open(); err != nil {
		log.Fatal(err)
		return
	}

	tasks := model.ResponseTasks{Tasks: service.GetTasks(s.Store.GetDB())}

	jsonResponse, err := json.Marshal(tasks)
	if err != nil {
		log.Fatal(err)
		return
	}

	if _, err = w.Write(jsonResponse); err != nil {
		log.Fatal(err)
		return
	}
}

// func handleGetGroups handles POST request for /groups
func (s *APIServer) handlePostGroups(w http.ResponseWriter, req *http.Request) {
	if err := s.Store.Open(); err != nil {
		log.Fatal(err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var group model.Groups
	err := decoder.Decode(&group)
	if err != nil {
		http.Error(w, "Wrong request struct", http.StatusBadRequest)
		log.Println(err)
		return
	}

	if group.Title == "" {
		http.Error(w, "title required", http.StatusBadRequest)
		return
	}

	group, err = service.PostGroups(s.Store.GetDB(), group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse, err := json.Marshal(group)
	if err != nil {
		log.Fatal(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, err = w.Write(jsonResponse); err != nil {
		log.Fatal(err)
		return
	}
}

// func handleGetGroups handles POST request for /tasks
func (s *APIServer) handlePostTasks(w http.ResponseWriter, req *http.Request) {
	if err := s.Store.Open(); err != nil {
		log.Fatal(err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var task model.Tasks
	err := decoder.Decode(&task)
	if err != nil {
		http.Error(w, "Wrong request struct", http.StatusBadRequest)
		log.Println(err)
		return
	}

	if task.Title == "" {
		http.Error(w, "title required", http.StatusBadRequest)
		return
	}

	task, err = service.PostTasks(s.Store.GetDB(), task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	jsonResponse, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, err = w.Write(jsonResponse); err != nil {
		log.Fatal(err)
		return
	}
}

// func handleGetGroups handles POST request for /time_frames
func (s *APIServer) handlePostTimeFrames(w http.ResponseWriter, req *http.Request) {
	if err := s.Store.Open(); err != nil {
		log.Fatal(err)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var timeFrame model.TimeFrames
	err := decoder.Decode(&timeFrame)
	if err != nil {
		http.Error(w, "Wrong time format", http.StatusBadRequest)
		return
	}

	if timeFrame.TaskId == "" {
		http.Error(w, "task id required", http.StatusBadRequest)
		return
	}

	if timeFrame.To.IsZero() || timeFrame.From.IsZero() {
		http.Error(w, "'From time' and 'To time' should be identified", http.StatusBadRequest)
		return
	}

	if timeFrame.From.After(timeFrame.To) {
		http.Error(w, "'From time' is after 'To time'", http.StatusBadRequest)
		return
	}

	timeFrame, err = service.PostTimeFrames(s.Store.GetDB(), timeFrame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	jsonResponse, err := json.Marshal(timeFrame)
	if err != nil {
		log.Fatal(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, err = w.Write(jsonResponse); err != nil {
		log.Fatal(err)
		return
	}
}

// func handleGetGroups handles PUT request for /groups/{id}
func (s *APIServer) handlePutGroups(w http.ResponseWriter, req *http.Request) {
	if err := s.Store.Open(); err != nil {
		log.Fatal(err)
		return
	}

	id := getID(req, "id")

	decoder := json.NewDecoder(req.Body)
	var group model.Groups
	err := decoder.Decode(&group)
	if err != nil {
		http.Error(w, "wrong request struct", http.StatusBadRequest)
		return
	}
	if group.Title == "" {
		http.Error(w, "title required", http.StatusBadRequest)
		return
	}

	group, err = service.PutGroups(s.Store.GetDB(), id, group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse, err := json.Marshal(group)
	if err != nil {
		log.Fatal(err)
		return
	}

	if _, err = w.Write(jsonResponse); err != nil {
		log.Fatal(err)
		return
	}

}

// func handleGetGroups handles PUT request for /tasks/{id}
func (s *APIServer) handlePutTasks(w http.ResponseWriter, req *http.Request) {
	if err := s.Store.Open(); err != nil {
		log.Fatal(err)
		return
	}

	id := getID(req, "id")

	decoder := json.NewDecoder(req.Body)
	var task model.Tasks
	err := decoder.Decode(&task)
	if err != nil {
		http.Error(w, "Wrong request struct", http.StatusBadRequest)
		log.Println(err)
		return
	}

	if task.Title == "" {
		http.Error(w, "title required", http.StatusBadRequest)
		return
	}

	if task.GroupId == "" {
		http.Error(w, "group id required", http.StatusBadRequest)
		return
	}

	task, err = service.PutTasks(s.Store.GetDB(), id, task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponse, err := json.Marshal(task)
	if err != nil {
		log.Fatal(err)
	}

	if _, err = w.Write(jsonResponse); err != nil {
		log.Fatal(err)
	}

}

// func handleGetGroups handles DELETE request for /groups/{id}
func (s *APIServer) handleDeleteGroups(w http.ResponseWriter, req *http.Request) {
	if err := s.Store.Open(); err != nil {
		log.Fatal(err)
	}

	id := getID(req, "id")

	if err := service.DeleteGroups(s.Store.GetDB(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// func handleGetGroups handles DELETE request for /tasks/{id}
func (s *APIServer) handleDeleteTasks(w http.ResponseWriter, req *http.Request) {
	if err := s.Store.Open(); err != nil {
		log.Fatal(err)
	}

	id := getID(req, "id")

	if err := service.DeleteTasks(s.Store.GetDB(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// func handleGetGroups handles DELETE request for /timeframes/{id}
func (s *APIServer) handleDeleteTimeFrames(w http.ResponseWriter, req *http.Request) {
	if err := s.Store.Open(); err != nil {
		log.Fatal(err)
		return
	}

	id := getID(req, "id")

	if err := service.DeleteTimeFrames(s.Store.GetDB(), id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

// func getID returns id of an object from url
func getID(req *http.Request, idName string) string {
	vars := mux.Vars(req)
	id := vars[idName]
	return id
}
