package handlrs

import (
	"net/http"
)

type TaskHandler struct {
}

func (h *TaskHandler) List() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "No Implemented", http.StatusNotImplemented)
	})
}

func (h *TaskHandler) Create() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "No Implemented", http.StatusNotImplemented)
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
