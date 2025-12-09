package Repositories

import (
	"errors"
	"sync"
	"task-manager/Domain"
)

type TaskRepository interface {
	Create(t Domain.Task) (Domain.Task, error)
	GetAll() []Domain.Task
	GetByID(id int) (Domain.Task, error)
	Update(t Domain.Task) error
	Delete(id int) error
}

type inMemoryTaskRepository struct {
	mu    sync.RWMutex
	tasks []Domain.Task
	last  int
}

func NewInMemoryTaskRepository() TaskRepository {
	return &inMemoryTaskRepository{
		tasks: make([]Domain.Task, 0),
		last:  0,
	}
}

func (r *inMemoryTaskRepository) Create(t Domain.Task) (Domain.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.last++
	t.ID = r.last
	r.tasks = append(r.tasks, t)
	return t, nil
}

func (r *inMemoryTaskRepository) GetAll() []Domain.Task {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Domain.Task, len(r.tasks))
	copy(out, r.tasks)
	return out
}

func (r *inMemoryTaskRepository) GetByID(id int) (Domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, t := range r.tasks {
		if t.ID == id {
			return t, nil
		}
	}
	return Domain.Task{}, errors.New("not found")
}

func (r *inMemoryTaskRepository) Update(t Domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.tasks {
		if r.tasks[i].ID == t.ID {
			r.tasks[i].Title = t.Title
			r.tasks[i].Description = t.Description
			return nil
		}
	}
	return errors.New("not found")
}

func (r *inMemoryTaskRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.tasks {
		if r.tasks[i].ID == id {
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}
