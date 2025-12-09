package Repositories

import (
	"errors"
	"sync"
	"task-manager/Domain"
)

type UserRepository interface {
	Create(u Domain.User) (Domain.User, error)
	FindByUsername(username string) (Domain.User, error)
	FindByID(id int) (Domain.User, error)
	Promote(id int) error
	GetAll() []Domain.User
}

type inMemoryUserRepository struct {
	mu    sync.RWMutex
	users []Domain.User
	last  int
}

func NewInMemoryUserRepository() UserRepository {
	return &inMemoryUserRepository{
		users: make([]Domain.User, 0),
		last:  0,
	}
}

func (r *inMemoryUserRepository) Create(u Domain.User) (Domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// unique username
	for _, x := range r.users {
		if x.Username == u.Username {
			return Domain.User{}, errors.New("username already exists")
		}
	}

	r.last++
	u.ID = r.last
	r.users = append(r.users, u)
	return u, nil
}

func (r *inMemoryUserRepository) FindByUsername(username string) (Domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.Username == username {
			return u, nil
		}
	}
	return Domain.User{}, errors.New("not found")
}

func (r *inMemoryUserRepository) FindByID(id int) (Domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.ID == id {
			return u, nil
		}
	}
	return Domain.User{}, errors.New("not found")
}

func (r *inMemoryUserRepository) Promote(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range r.users {
		if r.users[i].ID == id {
			r.users[i].Role = "admin"
			return nil
		}
	}
	return errors.New("not found")
}

func (r *inMemoryUserRepository) GetAll() []Domain.User {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Domain.User, len(r.users))
	copy(out, r.users)
	return out
}
