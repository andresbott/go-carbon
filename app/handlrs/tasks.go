package handlrs

import (
	"encoding/json"
	"fmt"
	"git.andresbott.com/Golang/carbon/internal/tasks"
	"git.andresbott.com/Golang/carbon/libs/auth"
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
	Id string `json:"id"`
}

func (h *TaskHandler) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		uData, err := auth.CtxUserInfo(r)
		if err != nil {
			http.Error(w, fmt.Sprintf("unable to get user data from request context: %s", err), http.StatusInternalServerError)
			return
		}

		if uData.UserId == "" && uData.IsAuthenticated {
			http.Error(w, "user information not in context", http.StatusBadRequest)
			return
		}

		payload := localTaskInput{}
		err = json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if payload.Text == "" {
			http.Error(w, "text cannot be empty in task payload", http.StatusBadRequest)
			return
		}
		t := tasks.Task{
			Text: payload.Text,
			Done: payload.Done,
		}
		Id, err := h.TaskManager.Create(&t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		output := localTaskOutput{
			Id: Id,
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
		http.Error(w, "No Implemented", http.StatusNotImplemented)
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
