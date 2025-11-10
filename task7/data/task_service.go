package data

import (
	"sync"
	"github.com/hababisha/task_manager/models"
	"strconv"
)

type TaskService struct {
	mu sync.Mutex
	tasks map[int]models.Task
	nextID int
}

func NewTaskService() *TaskService {
	return &TaskService{
		tasks: make(map[int]models.Task),
		nextID: 1,
	}
}

func (s *TaskService) GetAllTasks() []models.Task{
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks := make([]models.Task, 0, len(s.tasks))

	for _, task := range s.tasks{
		tasks = append(tasks, task)
	}

	return tasks
}

func (s *TaskService) GetTaskByID(id int) (models.Task, bool){
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[id]
	return task, exists
}

func (s *TaskService) CreateTask(task models.Task) models.Task{
	s.mu.Lock()
	defer s.mu.Unlock()

	task.ID = s.nextID
	s.nextID++
	s.tasks[task.ID] = task
	return task
}

func (s *TaskService) UpdateTask(id int, updated models.Task) (models.Task, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	task, exists := s.tasks[id]
	if !exists{
		return models.Task{}, false
	}

	if updated.Title != ""{
		task.Title = updated.Title
	}

	if updated.Description != ""{
		task.Description = updated.Description
	}

	if updated.Status != ""{
		task.Status = updated.Status
	}

	s.tasks[id] = task
	return task, true
}


//Delete Task
func (s *TaskService) DeleteTask(id int) bool{
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.tasks[id]; !exists {
		return false
	}
	delete(s.tasks, id)
	return true
}