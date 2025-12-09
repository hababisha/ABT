package Usecases

import (
	"errors"
	"task-manager/Domain"
)

type TaskUsecase interface {
	Create(title, description string) (Domain.Task, error)
	GetAll() ([]Domain.Task, error)
	GetByID(id int) (Domain.Task, error)
	Update(id int, title, description string) error
	Delete(id int) error
}

type taskUsecase struct {
	repo TaskRepo
}

type TaskRepo interface {
	Create(t Domain.Task) (Domain.Task, error)
	GetAll() []Domain.Task
	GetByID(id int) (Domain.Task, error)
	Update(t Domain.Task) error
	Delete(id int) error
}

func NewTaskUsecase(r TaskRepo) TaskUsecase {
	return &taskUsecase{repo: r}
}

func (t *taskUsecase) Create(title, description string) (Domain.Task, error) {
	task := Domain.Task{
		Title:       title,
		Description: description,
	}
	return t.repo.Create(task)
}

func (t *taskUsecase) GetAll() ([]Domain.Task, error) {
	return t.repo.GetAll(), nil
}

func (t *taskUsecase) GetByID(id int) (Domain.Task, error) {
	task, err := t.repo.GetByID(id)
	if err != nil {
		return Domain.Task{}, errors.New("not found")
	}
	return task, nil
}

func (t *taskUsecase) Update(id int, title, description string) error {
	_, err := t.repo.GetByID(id)
	if err != nil {
		return errors.New("not found")
	}
	return t.repo.Update(Domain.Task{ID: id, Title: title, Description: description})
}

func (t *taskUsecase) Delete(id int) error {
	_, err := t.repo.GetByID(id)
	if err != nil {
		return errors.New("not found")
	}
	return t.repo.Delete(id)
}
