package tasks_test

import (
	"errors"
	"fmt"
	"git.andresbott.com/Golang/carbon/internal/model/tasks"
	"git.andresbott.com/Golang/carbon/libs/logzero"
	"github.com/google/go-cmp/cmp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"testing"
)

const inMemorySqlite = "file::memory:"

func TestListTasks(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(inMemorySqlite), &gorm.Config{
		Logger: logzero.NewZeroGorm(*logzero.DefaultLogger(logzero.InfoLevel, nil), logzero.Cfg{IgnoreRecordNotFoundError: true}),
	})
	if err != nil {
		t.Fatal(err)
	}
	mngr, err := tasks.New(db, &sync.Mutex{})
	if err != nil {
		t.Fatal(err)
	}

	const User1 = "u1"
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 20; i++ {
			_ = createTask(t, mngr, "task"+strconv.Itoa(i)+"_"+User1, User1)
		}
	}()

	const User2 = "u2"
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 30; i++ {
			_ = createTask(t, mngr, "task"+strconv.Itoa(i)+"_"+User2, User2)
		}
	}()
	wg.Wait()

	t.Run("head of index", func(t *testing.T) {
		items, err := mngr.List(User1, 2, 0)
		if err != nil {
			t.Fatal(err)
		}
		got := []string{}
		for _, item := range items {
			got = append(got, item.Text)
		}

		want := []string{
			"task1_u1",
			"task2_u1",
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}
	})

	t.Run("differnt user and page", func(t *testing.T) {
		items, err := mngr.List(User2, 3, 1)
		if err != nil {
			t.Fatal(err)
		}
		got := []string{}
		for _, item := range items {
			got = append(got, item.Text)
		}

		want := []string{
			"task4_u2",
			"task5_u2",
			"task6_u2",
		}
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("unexpected value (-got +want)\n%s", diff)
		}
	})

}
func TestCrudTask(t *testing.T) {

	// initialize DB
	db, err := gorm.Open(sqlite.Open(inMemorySqlite), &gorm.Config{
		Logger: logzero.NewZeroGorm(*logzero.DefaultLogger(logzero.InfoLevel, nil), logzero.Cfg{IgnoreRecordNotFoundError: true}),
	})
	if err != nil {
		t.Fatal(err)
	}

	mngr, err := tasks.New(db, nil) // since in this test we run everything sequentially we explicitly don't set a mutex
	if err != nil {
		t.Fatal(err)
	}
	// create a bunch of tasks for different users
	t1 := createTask(t, mngr, "task1", "u1")
	t2 := createTask(t, mngr, "task2", "u1")
	t3 := createTask(t, mngr, "task1", "u2")

	// verify we can read the tasks
	readTask(t, mngr, t1, "u1", "task1", "")
	// make sure ownership is conserved
	readTask(t, mngr, t1, "u2", "", fmt.Sprintf("task with id: %s and owner u2 not found", t1))
	readTask(t, mngr, t2, "u1", "task2", "")
	readTask(t, mngr, t3, "u2", "task1", "")
	readTask(t, mngr, t3, "u3", "", fmt.Sprintf("task with id: %s and owner u3 not found", t3))

	// Complete a task
	setDone(t, mngr, t1, "u1", true, "")
	// send complete again
	setDone(t, mngr, t1, "u1", true, "")
	// wrong owner
	setDone(t, mngr, t1, "u2", true, fmt.Sprintf("task with id: %s and owner u2 not found", t1))
	// set to pending
	setDone(t, mngr, t1, "u1", false, "")
	// again
	setDone(t, mngr, t1, "u1", false, "")

	// update the text
	setText(t, mngr, t1, "u1", "task1Updated", "")
	setText(t, mngr, t1, "u2", "", fmt.Sprintf("task with id: %s and owner u2 not found", t1))

	// delete the task
	deleteTask(t, mngr, t1, "u2", fmt.Sprintf("task with id: %s and owner u2 not found", t1))
	deleteTask(t, mngr, t1, "u1", "")

}

func createTask(t *testing.T, mngr *tasks.Manager, content, owner string) string {
	task := tasks.Task{
		Text:    content,
		OwnerId: owner,
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

func readTask(t *testing.T, mngr *tasks.Manager, taskId, owner, want string, wantErr string) {
	// verify we can read the tasks
	task, err := mngr.Get(taskId, owner)
	if err != nil {
		if wantErr != err.Error() {
			t.Errorf("wanted error:\"%s\", but got: \"%s\"", wantErr, err.Error())
		}
	}
	if task.Text != want {
		t.Errorf("expect task value to be \"%s\", but got: \"%s\"", want, task.Text)
	}
}

func setDone(t *testing.T, mngr *tasks.Manager, taskId, owner string, val bool, wantErr string) {
	var err error
	if val {
		err = mngr.Complete(taskId, owner)
	} else {
		err = mngr.Pending(taskId, owner)
	}
	if err != nil {
		if wantErr != err.Error() {
			t.Errorf("wanted error:\"%s\", but got: \"%s\"", wantErr, err.Error())
		}
		return
	}
	task, err := mngr.Get(taskId, owner)
	if err != nil {
		t.Error(err)
	}
	if task.Done != val {
		t.Errorf("expect task value to be \"%t\", but got: \"%t\"", val, task.Done)
	}
}

func setText(t *testing.T, mngr *tasks.Manager, taskId, owner, text, wantErr string) {
	err := mngr.UpdateText(taskId, owner, text)
	if err != nil {
		if wantErr != err.Error() {
			t.Errorf("wanted error:\"%s\", but got: \"%s\"", wantErr, err.Error())
		}
		return
	}
	task, err := mngr.Get(taskId, owner)
	if err != nil {
		t.Error(err)
	}
	if task.Text != text {
		t.Errorf("expect task value to be \"%s\", but got: \"%s\"", text, task.Text)
	}
}

func deleteTask(t *testing.T, mngr *tasks.Manager, taskId, owner, wantErr string) {
	err := mngr.Delete(taskId, owner)
	if err != nil {
		if wantErr != err.Error() {
			t.Errorf("wanted error:\"%s\", but got: \"%s\"", wantErr, err.Error())
		}
		return
	}
	_, err = mngr.Get(taskId, owner)
	if err != nil {
		target := &tasks.TaskNotFound{}
		if !errors.As(err, &target) {
			t.Errorf("unexpected error: %v", err)
		}

	}

}
