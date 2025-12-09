package Usecases

import (
	"errors"
	"task-manager/Domain"
	"task-manager/Infrastructure"
)

type UserUsecase interface {
	Register(username, password string) (Domain.User, error)
	Login(username, password string) (string, error)
	Promote(userID int) error
}

type userUsecase struct {
	repo       UserRepo
	password   Infrastructure.PasswordService
	jwtService Infrastructure.JWTService
}

type UserRepo interface {
	Create(u Domain.User) (Domain.User, error)
	FindByUsername(username string) (Domain.User, error)
	FindByID(id int) (Domain.User, error)
	Promote(id int) error
	GetAll() []Domain.User
}

func NewUserUsecase(r UserRepo, p Infrastructure.PasswordService, j Infrastructure.JWTService) UserUsecase {
	return &userUsecase{repo: r, password: p, jwtService: j}
}

func (u *userUsecase) Register(username, password string) (Domain.User, error) {
	// first user becomes admin
	all := u.repo.GetAll()
	role := "user"
	if len(all) == 0 {
		role = "admin"
	}
	hashed, err := u.password.Hash(password)
	if err != nil {
		return Domain.User{}, err
	}
	user := Domain.User{
		Username: username,
		Password: hashed,
		Role:     role,
	}
	created, err := u.repo.Create(user)
	if err != nil {
		return Domain.User{}, err
	}
	// redact password
	created.Password = ""
	return created, nil
}

func (u *userUsecase) Login(username, password string) (string, error) {
	user, err := u.repo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}
	if !u.password.Compare(user.Password, password) {
		return "", errors.New("invalid credentials")
	}
	// generate JWT
	token, err := u.jwtService.Generate(user.ID, user.Username, user.Role)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *userUsecase) Promote(userID int) error {
	return u.repo.Promote(userID)
}
