package service

import (
	"dbmsbackend/domain"
	"dbmsbackend/repository"
	"dbmsbackend/util"
	"fmt"
)

type UserRepository interface {
	Create(*domain.User) (*domain.User, error)
	FetchByID(int) (*domain.User, error)
	FetchByEmail(string) (*domain.User, error)
	Update(int, string, string) (*domain.User, error)
	Query(map[string]interface{}) []domain.User
	Initialize(*util.Config) error
}

type User struct {
	repo UserRepository
}

func (userService *User) Initialize(config *util.Config) (err error) {

	if config.DBDriver == "sqlite" {

		userService.repo = new(repository.UserSqldbRepository)
		err = userService.repo.Initialize(config)

		if err != nil {
			err = fmt.Errorf("initializing user service: %w", err)
		}
	}

	return

}

func (userService *User) New(name string, email string, phone string, kind domain.UserType, password string) (*domain.User, error) {

	entity := domain.NewUser(name, email, phone, kind, password)

	entity, err := userService.repo.Create(entity)

	if err != nil {
		return nil, err
	}

	return entity, nil
}

func (userService *User) Update(id int, name string, phone string) (*domain.User, error) {

	entity, err := userService.repo.Update(id, name, phone)

	if err != nil {
		return nil, err
	}

	return entity, err
}

func (userService *User) Query(condition map[string]interface{}) []domain.User {

	users := userService.repo.Query(condition)

	return users
}

func (userService *User) GetByID(id int) (*domain.User, error) {
	entity, err := userService.repo.FetchByID(id)

	if err != nil {
		return entity, &util.ErrNotFound{
			Err: err,
		}
	}

	return entity, err
}

func (userService *User) Login(email string, password string) (*domain.User, error) {
	entity, err := userService.repo.FetchByEmail(email)

	if err != nil {
		return entity, err
	}

	if entity.VerifyPassword(password) {
		return entity, nil
	}

	return entity, err
}
