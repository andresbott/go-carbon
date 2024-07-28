package tasks

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Manager struct {
	db        *gorm.DB
	writeLock sync.Locker // sqlite does not allow concurrent write
}

func New(db *gorm.DB, writeLock sync.Locker) (*Manager, error) {
	// Migrate the schema
	err := db.AutoMigrate(&Task{})
	if err != nil {
		return nil, err
	}

	m := Manager{
		db:        db,
		writeLock: writeLock,
	}
	return &m, nil
}

type Task struct {
	ID      string `gorm:"primaryKey,index"`
	OwnerId string `gorm:"index"`
	Text    string
	Done    bool

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (user *Task) BeforeCreate(db *gorm.DB) (err error) {
	// UUID version 4
	user.ID = uuid.NewString()
	return
}

type TaskNotFound struct {
	id    string
	owner string
}

func (m *TaskNotFound) Error() string {
	return fmt.Sprintf("task with id: %s and owner %s not found", m.id, m.owner)
}

func (m Manager) List(owner string, size, page int) ([]Task, error) {
	if size <= 0 {
		size = 20
	}
	if size >= 50 {
		size = 50
	}

	offset := size * (page - 1)
	if offset <= 0 {
		offset = 0
	}
	tasks := make([]Task, size)
	result := m.db.Where("owner_id = ?", owner).Model(&Task{}).Offset(offset).Limit(size).Find(&tasks)
	if result.Error != nil {
		return nil, result.Error
	}
	return tasks, nil
}

func (m Manager) Create(task *Task) (string, error) {
	if m.writeLock != nil {
		m.writeLock.Lock()
		defer m.writeLock.Unlock()
	}

	result := m.db.Create(task)
	if result.Error != nil {
		return "", result.Error
	}
	return task.ID, nil
}

func (m Manager) Get(id, owner string) (Task, error) {
	t := Task{}
	result := m.db.First(&t, "ID = ? AND owner_id = ?", id, owner)
	if result.RowsAffected == 0 {
		return t, &TaskNotFound{id: id, owner: owner}
	}
	return t, nil
}

func (m Manager) Update(id, owner, text string, done *bool) error {
	if m.writeLock != nil {
		m.writeLock.Lock()
		defer m.writeLock.Unlock()
	}

	fieldMap := map[string]any{}
	if text != "" {
		fieldMap["text"] = text
	}
	if done != nil {
		fieldMap["done"] = *done
	}

	t := Task{}
	result := m.db.Model(&t).
		Where("ID = ? AND owner_id = ?", id, owner).
		Updates(fieldMap)

	if result.RowsAffected == 0 {
		return &TaskNotFound{id: id, owner: owner}
	}
	return nil
}

func (m Manager) Delete(id, owner string) error {
	if m.writeLock != nil {
		m.writeLock.Lock()
		defer m.writeLock.Unlock()
	}

	t := Task{}
	result := m.db.Where("ID = ? AND owner_id = ?", id, owner).Delete(&t)

	if result.RowsAffected == 0 {
		return &TaskNotFound{id: id, owner: owner}
	}
	return nil
}

// TODO, hard delete
