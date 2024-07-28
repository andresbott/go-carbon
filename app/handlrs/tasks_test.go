package handlrs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"git.andresbott.com/Golang/carbon/internal/model/tasks"
	"git.andresbott.com/Golang/carbon/libs/auth"
	"git.andresbott.com/Golang/carbon/libs/logzero"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
)

const (
	user1 = "user1"
	user2 = "user2"
)

func TestTaskHandler_List(t *testing.T) {
	tcs := []struct {
		name       string
		req        func() (*http.Request, error)
		expecErr   string
		expectCode int
		expect     localTaskList
	}{
		{
			name: "successful request",
			req: func() (*http.Request, error) {

				req, err := http.NewRequest("GET", "/api/tasks", nil)
				if err != nil {
					return nil, err
				}
				ctx := req.Context()
				ctx = context.WithValue(ctx, auth.UserIdKey, user1)
				ctx = context.WithValue(ctx, auth.UserIsLoggedInKey, true)
				req = req.WithContext(ctx)

				q := req.URL.Query()
				q.Add("limit", "2")
				req.URL.RawQuery = q.Encode()

				return req, nil
			},
			expectCode: http.StatusOK,
			expect: localTaskList{
				Count: 2,
				Tasks: []localTaskOutput{
					{Text: "task1_user1"},
					{Text: "task2_user1"},
				},
			},
		},
		{
			name: "success with different limit and page",
			req: func() (*http.Request, error) {

				req, err := http.NewRequest("GET", "/api/tasks", nil)
				if err != nil {
					return nil, err
				}
				ctx := req.Context()
				ctx = context.WithValue(ctx, auth.UserIdKey, user2)
				ctx = context.WithValue(ctx, auth.UserIsLoggedInKey, true)
				req = req.WithContext(ctx)

				q := req.URL.Query()
				q.Add(limitParam, "3")
				q.Add(pageParam, "2")
				req.URL.RawQuery = q.Encode()

				return req, nil
			},
			expectCode: http.StatusOK,
			expect: localTaskList{
				Count: 3,
				Tasks: []localTaskOutput{
					{Text: "task4_user2"},
					{Text: "task5_user2"},
					{Text: "task6_user2"},
				},
			},
		},
		{
			name: "not authenticated",
			req: func() (*http.Request, error) {

				req, err := http.NewRequest("GET", "/api/tasks", nil)
				if err != nil {
					return nil, err
				}
				ctx := req.Context()
				ctx = context.WithValue(ctx, auth.UserIdKey, user1)
				ctx = context.WithValue(ctx, auth.UserIsLoggedInKey, false)
				req = req.WithContext(ctx)

				q := req.URL.Query()
				q.Add("limit", "2")
				req.URL.RawQuery = q.Encode()

				return req, nil
			},
			expecErr:   "user login information not provided in request context: isLoggedIn",
			expectCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.req()
			if err != nil {
				t.Fatal(err)
			}
			th, err := taskHandler()
			if err != nil {
				t.Fatal(err)
			}

			for i := 1; i <= 30; i++ {
				_ = createTask(t, th.TaskManager, "task"+strconv.Itoa(i)+"_"+user1, user1)
			}
			for i := 1; i <= 30; i++ {
				_ = createTask(t, th.TaskManager, "task"+strconv.Itoa(i)+"_"+user2, user2)
			}

			recorder := httptest.NewRecorder()

			handler := th.List()
			handler.ServeHTTP(recorder, req)

			if tc.expecErr != "" {
				if status := recorder.Code; status != tc.expectCode {
					t.Errorf("handler returned wrong status code: got %v want %v",
						status, tc.expectCode)
				}
				respText, err := io.ReadAll(recorder.Body)
				if err != nil {
					t.Fatal(err)
				}
				got := strings.TrimSuffix(string(respText), "\n")
				if got != tc.expecErr {
					t.Errorf("unexpecter error message: got \"%s\" want \"%v\"",
						got, tc.expecErr)
				}

			} else {

				if status := recorder.Code; status != tc.expectCode {
					t.Errorf("handler returned wrong status code: got %v want %v",
						status, tc.expectCode)
				}

				got := localTaskList{}
				err = json.NewDecoder(recorder.Body).Decode(&got)
				if err != nil {
					t.Fatal(err)
				}
				if diff := cmp.Diff(got, tc.expect, cmpopts.IgnoreFields(localTaskOutput{}, "Id")); diff != "" {
					t.Errorf("unexpected value (-got +want)\n%s", diff)
				}

			}

		})
	}
}

func TestTaskHandler_Create(t *testing.T) {
	tcs := []struct {
		name       string
		req        func() (*http.Request, error)
		expecErr   string
		expectCode int
	}{
		{
			name: "successful request",
			req: func() (*http.Request, error) {
				var jsonStr = []byte(`{"text":"some task"}`)
				req, err := http.NewRequest("PUT", "/api/tasks", bytes.NewBuffer(jsonStr))
				if err != nil {
					return nil, err
				}
				ctx := req.Context()
				ctx = context.WithValue(ctx, auth.UserIdKey, "user1")
				ctx = context.WithValue(ctx, auth.UserIsLoggedInKey, true)
				req = req.WithContext(ctx)
				return req, nil
			},
			expectCode: http.StatusOK,
		},
		{
			name: "not authenticated",
			req: func() (*http.Request, error) {
				// this scenario the middle ware should have returned already a 401
				var jsonStr = []byte(`{"text":"some task"}`)
				req, err := http.NewRequest("PUT", "/api/tasks", bytes.NewBuffer(jsonStr))
				if err != nil {
					return nil, err
				}
				ctx := req.Context()
				ctx = context.WithValue(ctx, auth.UserIsLoggedInKey, false)
				req = req.WithContext(ctx)
				return req, nil
			},
			expecErr:   "user login information not provided in request context: isLoggedIn",
			expectCode: http.StatusInternalServerError,
		},
		{
			name: "empty payload",
			req: func() (*http.Request, error) {

				req, err := http.NewRequest("PUT", "/api/tasks", nil)
				if err != nil {
					return nil, err
				}
				ctx := req.Context()
				ctx = context.WithValue(ctx, auth.UserIdKey, "user1")
				ctx = context.WithValue(ctx, auth.UserIsLoggedInKey, true)
				req = req.WithContext(ctx)
				return req, nil
			},
			expecErr:   "request had empty body",
			expectCode: http.StatusBadRequest,
		},
		{
			name: "malformed payload",
			req: func() (*http.Request, error) {
				var jsonStr = []byte(`{`)
				req, err := http.NewRequest("PUT", "/api/tasks", bytes.NewBuffer(jsonStr))
				if err != nil {
					return nil, err
				}
				ctx := req.Context()
				ctx = context.WithValue(ctx, auth.UserIdKey, "user1")
				ctx = context.WithValue(ctx, auth.UserIsLoggedInKey, true)
				req = req.WithContext(ctx)
				return req, nil
			},
			expecErr:   "unable to decode json: unexpected EOF",
			expectCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			req, err := tc.req()
			if err != nil {
				t.Fatal(err)
			}

			th, err := taskHandler()
			if err != nil {
				t.Fatal(err)
			}
			recorder := httptest.NewRecorder()

			handler := th.Create()
			handler.ServeHTTP(recorder, req)

			if tc.expecErr != "" {
				if status := recorder.Code; status != tc.expectCode {
					t.Errorf("handler returned wrong status code: got %v want %v",
						status, tc.expectCode)
				}
				respText, err := io.ReadAll(recorder.Body)
				if err != nil {
					t.Fatal(err)
				}
				got := strings.TrimSuffix(string(respText), "\n")
				if got != tc.expecErr {
					t.Errorf("unexpecter error message: got \"%s\" want \"%v\"",
						got, tc.expecErr)
				}

			} else {

				if status := recorder.Code; status != tc.expectCode {
					t.Errorf("handler returned wrong status code: got %v want %v",
						status, tc.expectCode)
				}

				task := tasks.Task{}
				err = json.NewDecoder(recorder.Body).Decode(&task)
				if err != nil {
					t.Fatal(err)
				}
				if !IsValidUUID(task.ID) {
					t.Error("returned task ID is not a valid UUID")
				}

			}

		})
	}
}

func TestTaskHandler_Read(t *testing.T) {
	tcs := []struct {
		name       string
		req        func(id string) (*http.Request, error)
		expecErr   string
		expectCode int
	}{
		{
			name: "successful request",
			req: func(id string) (*http.Request, error) {

				req, err := http.NewRequest("GET", "/api/tasks/"+id, nil)
				if err != nil {
					return nil, err
				}
				ctx := req.Context()
				ctx = context.WithValue(ctx, auth.UserIdKey, "user1")
				ctx = context.WithValue(ctx, auth.UserIsLoggedInKey, true)
				req = req.WithContext(ctx)
				req = mux.SetURLVars(req, map[string]string{
					"ID": id,
				})
				return req, nil
			},
			expectCode: http.StatusOK,
		},
		{
			name: "fail for other user",
			req: func(id string) (*http.Request, error) {

				req, err := http.NewRequest("GET", "/api/tasks/"+id, nil)
				if err != nil {
					return nil, err
				}
				ctx := req.Context()
				ctx = context.WithValue(ctx, auth.UserIdKey, "user2")
				ctx = context.WithValue(ctx, auth.UserIsLoggedInKey, true)
				req = req.WithContext(ctx)
				req = mux.SetURLVars(req, map[string]string{
					"ID": id,
				})
				return req, nil
			},
			expecErr:   "task with id: %s and owner user2 not found",
			expectCode: http.StatusNotFound,
		},
		{
			name: "not authenticated",
			req: func(id string) (*http.Request, error) {

				req, err := http.NewRequest("GET", "/api/tasks/"+id, nil)
				if err != nil {
					return nil, err
				}
				ctx := req.Context()
				ctx = context.WithValue(ctx, auth.UserIdKey, "user1")
				ctx = context.WithValue(ctx, auth.UserIsLoggedInKey, false)
				req = req.WithContext(ctx)
				req = mux.SetURLVars(req, map[string]string{
					"ID": id,
				})
				return req, nil
			},
			expecErr:   "user login information not provided in request context: isLoggedIn",
			expectCode: http.StatusInternalServerError,
		},
		{
			name: "empty task ID",
			req: func(id string) (*http.Request, error) {

				req, err := http.NewRequest("GET", "/api/tasks/"+id, nil)
				if err != nil {
					return nil, err
				}
				ctx := req.Context()
				ctx = context.WithValue(ctx, auth.UserIdKey, "user1")
				ctx = context.WithValue(ctx, auth.UserIsLoggedInKey, true)
				req = req.WithContext(ctx)
				req = mux.SetURLVars(req, map[string]string{
					"ID": "",
				})
				return req, nil
			},
			expecErr:   "no task id provided",
			expectCode: http.StatusBadRequest,
		},
		{
			name: "malformed payload",
			req: func(id string) (*http.Request, error) {

				req, err := http.NewRequest("GET", "/api/tasks/"+id, nil)
				if err != nil {
					return nil, err
				}
				ctx := req.Context()
				ctx = context.WithValue(ctx, auth.UserIdKey, "user1")
				ctx = context.WithValue(ctx, auth.UserIsLoggedInKey, true)
				req = req.WithContext(ctx)
				req = mux.SetURLVars(req, map[string]string{
					"ID": "ddd",
				})
				return req, nil
			},
			expecErr:   "task id is not a UUID",
			expectCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			th, err := taskHandler()
			if err != nil {
				t.Fatal(err)
			}
			sampleTask := tasks.Task{
				Text:    "sample",
				OwnerId: "user1",
			}
			taskId, err := th.TaskManager.Create(&sampleTask)
			if err != nil {
				t.Fatal(err)
			}

			req, err := tc.req(taskId)
			if err != nil {
				t.Fatal(err)
			}
			recorder := httptest.NewRecorder()

			handler := th.Read()
			handler.ServeHTTP(recorder, req)

			if tc.expecErr != "" {
				if status := recorder.Code; status != tc.expectCode {
					t.Errorf("handler returned wrong status code: got %v want %v",
						status, tc.expectCode)
				}
				respText, err := io.ReadAll(recorder.Body)
				if err != nil {
					t.Fatal(err)
				}
				got := strings.TrimSuffix(string(respText), "\n")

				want := tc.expecErr
				if strings.Contains(tc.expecErr, "%") {
					want = fmt.Sprintf(tc.expecErr, taskId)
				}
				if got != want {
					t.Errorf("unexpecter error message: got \"%s\" want \"%v\"",
						got, want)
				}

			} else {

				if status := recorder.Code; status != tc.expectCode {
					t.Errorf("handler returned wrong status code: got %v want %v",
						status, tc.expectCode)
				}

				got := tasks.Task{}
				err = json.NewDecoder(recorder.Body).Decode(&got)
				if err != nil {
					t.Fatal(err)
				}
				want := tasks.Task{
					ID:   taskId,
					Text: "sample",
				}
				if diff := cmp.Diff(got, want); diff != "" {
					t.Errorf("unexpected value (-got +want)\n%s", diff)
				}
			}

		})
	}
}

func TestTaskHandler_Update(t *testing.T) {
	tcs := []struct {
		name       string
		req        func(id string) (*http.Request, error)
		expecErr   string
		expectCode int
	}{
		{
			name: "successful request",
			req: func(id string) (*http.Request, error) {
				var jsonStr = []byte(`{"text":"updated text"}`)
				req, err := http.NewRequest("PUT", "/api/tasks/"+id, bytes.NewBuffer(jsonStr))
				if err != nil {
					return nil, err
				}
				ctx := req.Context()
				ctx = context.WithValue(ctx, auth.UserIdKey, "user1")
				ctx = context.WithValue(ctx, auth.UserIsLoggedInKey, true)
				req = req.WithContext(ctx)
				req = mux.SetURLVars(req, map[string]string{
					"ID": id,
				})
				return req, nil
			},
			expectCode: http.StatusAccepted,
		},
		{
			name: "not authenticated",
			req: func(id string) (*http.Request, error) {
				var jsonStr = []byte(`{"text":"updated text"}`)
				req, err := http.NewRequest("PUT", "/api/tasks/"+id, bytes.NewBuffer(jsonStr))
				if err != nil {
					return nil, err
				}
				ctx := req.Context()
				ctx = context.WithValue(ctx, auth.UserIdKey, "user1")
				ctx = context.WithValue(ctx, auth.UserIsLoggedInKey, false)
				req = req.WithContext(ctx)
				req = mux.SetURLVars(req, map[string]string{
					"ID": id,
				})
				return req, nil
			},
			expecErr:   "user login information not provided in request context: isLoggedIn",
			expectCode: http.StatusInternalServerError,
		},
		{
			name: "empty payload",
			req: func(id string) (*http.Request, error) {
				req, err := http.NewRequest("PUT", "/api/tasks/"+id, nil)
				if err != nil {
					return nil, err
				}
				ctx := req.Context()
				ctx = context.WithValue(ctx, auth.UserIdKey, "user1")
				ctx = context.WithValue(ctx, auth.UserIsLoggedInKey, true)
				req = req.WithContext(ctx)
				req = mux.SetURLVars(req, map[string]string{
					"ID": id,
				})
				return req, nil
			},
			expecErr:   "request had empty body",
			expectCode: http.StatusBadRequest,
		},
		{
			name: "malformed payload",
			req: func(id string) (*http.Request, error) {
				var jsonStr = []byte(`{"text":"updated te`)
				req, err := http.NewRequest("PUT", "/api/tasks/"+id, bytes.NewBuffer(jsonStr))
				if err != nil {
					return nil, err
				}
				ctx := req.Context()
				ctx = context.WithValue(ctx, auth.UserIdKey, "user1")
				ctx = context.WithValue(ctx, auth.UserIsLoggedInKey, true)
				req = req.WithContext(ctx)
				req = mux.SetURLVars(req, map[string]string{
					"ID": id,
				})
				return req, nil
			},
			expecErr:   "unable to decode json: unexpected EOF",
			expectCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			th, err := taskHandler()
			if err != nil {
				t.Fatal(err)
			}
			sampleTask := tasks.Task{
				Text:    "sample",
				OwnerId: "user1",
			}
			taskId, err := th.TaskManager.Create(&sampleTask)
			if err != nil {
				t.Fatal(err)
			}

			req, err := tc.req(taskId)
			if err != nil {
				t.Fatal(err)
			}
			recorder := httptest.NewRecorder()

			handler := th.Update()
			handler.ServeHTTP(recorder, req)

			if tc.expecErr != "" {
				if status := recorder.Code; status != tc.expectCode {
					t.Errorf("handler returned wrong status code: got %v want %v",
						status, tc.expectCode)
				}
				respText, err := io.ReadAll(recorder.Body)
				if err != nil {
					t.Fatal(err)
				}
				got := strings.TrimSuffix(string(respText), "\n")
				if got != tc.expecErr {
					t.Errorf("unexpecter error message: got \"%s\" want \"%v\"",
						got, tc.expecErr)
				}

			} else {

				if status := recorder.Code; status != tc.expectCode {
					t.Errorf("handler returned wrong status code: got %v want %v",
						status, tc.expectCode)
				}
			}

		})
	}
}

func TestTaskHandler_Delete(t *testing.T) {
	tcs := []struct {
		name       string
		req        func(id string) (*http.Request, error)
		expecErr   string
		expectCode int
	}{
		{
			name: "successful request",
			req: func(id string) (*http.Request, error) {

				req, err := http.NewRequest("GET", "/api/tasks/"+id, nil)
				if err != nil {
					return nil, err
				}
				ctx := req.Context()
				ctx = context.WithValue(ctx, auth.UserIdKey, "user1")
				ctx = context.WithValue(ctx, auth.UserIsLoggedInKey, true)
				req = req.WithContext(ctx)
				req = mux.SetURLVars(req, map[string]string{
					"ID": id,
				})
				return req, nil
			},
			expectCode: http.StatusAccepted,
		},
		{
			name: "fail for other user",
			req: func(id string) (*http.Request, error) {

				req, err := http.NewRequest("GET", "/api/tasks/"+id, nil)
				if err != nil {
					return nil, err
				}
				ctx := req.Context()
				ctx = context.WithValue(ctx, auth.UserIdKey, "user2")
				ctx = context.WithValue(ctx, auth.UserIsLoggedInKey, true)
				req = req.WithContext(ctx)
				req = mux.SetURLVars(req, map[string]string{
					"ID": id,
				})
				return req, nil
			},
			expecErr:   "task with id: %s and owner user2 not found",
			expectCode: http.StatusNotFound,
		},
		{
			name: "not authenticated",
			req: func(id string) (*http.Request, error) {

				req, err := http.NewRequest("GET", "/api/tasks/"+id, nil)
				if err != nil {
					return nil, err
				}
				ctx := req.Context()
				ctx = context.WithValue(ctx, auth.UserIdKey, "user1")
				ctx = context.WithValue(ctx, auth.UserIsLoggedInKey, false)
				req = req.WithContext(ctx)
				req = mux.SetURLVars(req, map[string]string{
					"ID": id,
				})
				return req, nil
			},
			expecErr:   "user login information not provided in request context: isLoggedIn",
			expectCode: http.StatusInternalServerError,
		},
		{
			name: "empty task ID",
			req: func(id string) (*http.Request, error) {

				req, err := http.NewRequest("GET", "/api/tasks/"+id, nil)
				if err != nil {
					return nil, err
				}
				ctx := req.Context()
				ctx = context.WithValue(ctx, auth.UserIdKey, "user1")
				ctx = context.WithValue(ctx, auth.UserIsLoggedInKey, true)
				req = req.WithContext(ctx)
				req = mux.SetURLVars(req, map[string]string{
					"ID": "",
				})
				return req, nil
			},
			expecErr:   "no task id provided",
			expectCode: http.StatusBadRequest,
		},
		{
			name: "malformed payload",
			req: func(id string) (*http.Request, error) {

				req, err := http.NewRequest("GET", "/api/tasks/"+id, nil)
				if err != nil {
					return nil, err
				}
				ctx := req.Context()
				ctx = context.WithValue(ctx, auth.UserIdKey, "user1")
				ctx = context.WithValue(ctx, auth.UserIsLoggedInKey, true)
				req = req.WithContext(ctx)
				req = mux.SetURLVars(req, map[string]string{
					"ID": "ddd",
				})
				return req, nil
			},
			expecErr:   "task id is not a UUID",
			expectCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			th, err := taskHandler()
			if err != nil {
				t.Fatal(err)
			}
			sampleTask := tasks.Task{
				Text:    "sample",
				OwnerId: "user1",
			}
			taskId, err := th.TaskManager.Create(&sampleTask)
			if err != nil {
				t.Fatal(err)
			}

			req, err := tc.req(taskId)
			if err != nil {
				t.Fatal(err)
			}
			recorder := httptest.NewRecorder()

			handler := th.Delete()
			handler.ServeHTTP(recorder, req)

			if tc.expecErr != "" {
				if status := recorder.Code; status != tc.expectCode {
					t.Errorf("handler returned wrong status code: got %v want %v",
						status, tc.expectCode)
				}
				respText, err := io.ReadAll(recorder.Body)
				if err != nil {
					t.Fatal(err)
				}
				got := strings.TrimSuffix(string(respText), "\n")

				want := tc.expecErr
				if strings.Contains(tc.expecErr, "%") {
					want = fmt.Sprintf(tc.expecErr, taskId)
				}
				if got != want {
					t.Errorf("unexpecter error message: got \"%s\" want \"%v\"",
						got, want)
				}

			} else {

				if status := recorder.Code; status != tc.expectCode {
					t.Errorf("handler returned wrong status code: got %v want %v",
						status, tc.expectCode)
				}
			}
		})
	}
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

const inMemorySqlite = "file::memory:"

func taskHandler() (*TaskHandler, error) {
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
	th := TaskHandler{
		TaskManager: mngr,
	}
	return &th, nil
}

func createTask(t *testing.T, mngr *tasks.Manager, content, owner string) string {
	b := false
	task := tasks.Task{
		Text:    content,
		OwnerId: owner,
		Done:    b,
	}
	id, err := mngr.Create(&task)
	if err != nil {
		t.Fatal(err)
	}
	if id == "" {
		t.Error("returned id should not be empty")
	}
	return id
}
