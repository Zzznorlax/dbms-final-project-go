package repository

import (
	"dbmsbackend/database/sqldb"
	"dbmsbackend/domain"
	"dbmsbackend/util"
	"fmt"

	"gorm.io/gorm"
)

type UserSqldbRepository struct {
	db *gorm.DB
}

func (repo *UserSqldbRepository) Initialize(config *util.Config) (err error) {
	dbGetter := sqldb.NewSQLiteGetter()
	repo.db = dbGetter(config.DBSource)

	if err != nil {
		err = fmt.Errorf("connecting to db: %w", err)
	}

	repo.db.AutoMigrate(&sqldb.User{})

	return err
}

func toUserDao(entity *domain.User) *sqldb.User {
	dao := sqldb.User{
		Name:           entity.Name,
		Email:          entity.Email,
		Phone:          entity.Phone,
		Kind:           string(entity.Kind),
		HashedPassword: entity.HashedPassword,
		Salt:           entity.Salt,
	}

	return &dao
}

func toUserEntity(dao *sqldb.User) *domain.User {
	entity := domain.User{
		ID:             int(dao.ID),
		Name:           dao.Name,
		Email:          dao.Email,
		Phone:          dao.Phone,
		Kind:           domain.UserType(dao.Kind),
		HashedPassword: dao.HashedPassword,
		Salt:           dao.Salt,
	}

	return &entity
}

func (repo *UserSqldbRepository) Create(user *domain.User) (*domain.User, error) {

	dao := toUserDao(user)

	if err := repo.db.Create(dao).Error; err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	return toUserEntity(dao), nil
}

func (repo *UserSqldbRepository) Update(id int, name string, phone string) (*domain.User, error) {

	dao := new(sqldb.User)

	if err := repo.db.First(&dao, id).Error; err != nil {
		return nil, fmt.Errorf("updating user: %w", err)
	}

	dao.Name = name
	dao.Phone = phone

	if err := repo.db.Save(&dao).Error; err != nil {
		return nil, fmt.Errorf("updating user: %w", err)
	}

	return toUserEntity(dao), nil
}

func (repo *UserSqldbRepository) Query(condition map[string]interface{}) []domain.User {

	var daos []sqldb.User
	var result []domain.User

	repo.db.Where(condition).Find(&daos)

	for _, item := range daos {
		result = append(result, *toUserEntity(&item))
	}

	return result
}

func (repo *UserSqldbRepository) FetchByID(id int) (*domain.User, error) {

	dao := new(sqldb.User)

	if err := repo.db.First(&dao, id).Error; err != nil {
		return nil, fmt.Errorf("fetching user %v: %w", id, err)
	}

	return toUserEntity(dao), nil
}

func (repo *UserSqldbRepository) FetchByEmail(email string) (*domain.User, error) {

	dao := new(sqldb.User)

	if err := repo.db.Where("email = ?", email).First(&dao).Error; err != nil {
		return nil, fmt.Errorf("fetching user %v: %w", email, err)
	}

	return toUserEntity(dao), nil
}
