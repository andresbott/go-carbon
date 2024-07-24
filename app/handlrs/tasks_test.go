package handlrs

import (
	"bytes"
	"context"
	"encoding/json"
	"git.andresbott.com/Golang/carbon/internal/tasks"
	"git.andresbott.com/Golang/carbon/libs/auth"
	"git.andresbott.com/Golang/carbon/libs/logzero"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestHealthCreateTask(t *testing.T) {

	// todo create mmore tests for scenarios:
	// not authenticated
	// check different errors

	var jsonStr = []byte(`{"text":"some task"}`)
	req, err := http.NewRequest("PUT", "/api/tasks", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Fatal(err)
	}
	ctx := req.Context()
	ctx = context.WithValue(ctx, auth.UserIdKey, "user1")
	ctx = context.WithValue(ctx, auth.UserIsLoggedInKey, true)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	tm, err := newTaskManager()
	if err != nil {
		t.Fatal(err)
	}
	th := TaskHandler{
		TaskManager: tm,
	}

	handler := th.Create()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	task := tasks.Task{}
	err = json.NewDecoder(rr.Body).Decode(&task)
	if err != nil {
		t.Fatal(err)
	}
	if task.ID == "" {
		t.Errorf("expect response to contain a task ID, but got empty sting")
	}

}

const inMemorySqlite = "file::memory:"

func newTaskManager() (*tasks.Manager, error) {
	db, err := gorm.Open(sqlite.Open(inMemorySqlite), &gorm.Config{
		Logger: logzero.NewZeroGorm(*logzero.DefaultLogger(logzero.InfoLevel, nil), logzero.Cfg{IgnoreRecordNotFoundError: true}),
	})
	if err != nil {
		return nil, err
	}
	mngr, err := tasks.New(db, &sync.Mutex{})
	if err != nil {
		return nil, err
	}
	return mngr, nil
}
