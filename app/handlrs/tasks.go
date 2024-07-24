package handlrs

import (
	"encoding/json"
	"errors"
	"fmt"
	"git.andresbott.com/Golang/carbon/internal/model/tasks"
	"git.andresbott.com/Golang/carbon/libs/auth"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
)

type TaskHandler struct {
	TaskManager *tasks.Manager
}

func (h *TaskHandler) List() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "No Implemented", http.StatusNotImplemented)
	})
}

type localTaskInput struct {
	Text string `json:"text"`
	Done bool
}
type localTaskOutput struct {
	Id   string `json:"id"`
	Text string `json:"text"`
	Done bool
}

func (h *TaskHandler) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uData, err := auth.CtxCheckAuth(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if r.Body == nil {
			http.Error(w, fmt.Sprintf("request had empty body"), http.StatusBadRequest)
			return
		}
		payload := localTaskInput{}
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to decode json: %s", err.Error()), http.StatusBadRequest)
			return
		}

		if payload.Text == "" {
			http.Error(w, "text cannot be empty req task payload", http.StatusBadRequest)
			return
		}
		t := tasks.Task{
			Text:    payload.Text,
			Done:    payload.Done,
			OwnerId: uData.UserId,
		}
		Id, err := h.TaskManager.Create(&t)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to store task in DB: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		output := localTaskOutput{
			Id:   Id,
			Text: payload.Text,
			Done: payload.Done,
		}
		respJson, err := json.Marshal(output)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write(respJson)
	})
}

func (h *TaskHandler) Read() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		taskId, ok := vars["ID"]
		if !ok {
			http.Error(w, fmt.Sprintf("could not extract id to read from request context"), http.StatusInternalServerError)
			return
		}
		if taskId == "" {
			http.Error(w, fmt.Sprint("no task id provided"), http.StatusBadRequest)
			return
		}
		_, err := uuid.Parse(taskId)
		if err != nil {
			http.Error(w, fmt.Sprint("task id is not a UUID"), http.StatusBadRequest)
			return
		}

		uData, err := auth.CtxCheckAuth(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		Task, err := h.TaskManager.Get(taskId, uData.UserId)
		if err != nil {
			t := &tasks.TaskNotFound{}
			if errors.As(err, &t) {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("unable to get Task: %s", err.Error()), http.StatusInternalServerError)
			}
			return
		}

		output := localTaskOutput{
			Id:   Task.ID,
			Text: Task.Text,
			Done: Task.Done,
		}
		respJson, err := json.Marshal(output)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respJson)
	})
}

func (h *TaskHandler) Update() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "No Implemented", http.StatusNotImplemented)
	})
}

func (h *TaskHandler) Delete() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "No Implemented", http.StatusNotImplemented)
	})
}
