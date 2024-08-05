package user

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	"errors"
)

type UserService struct {
	userRepo *UserRepo
}

//

// FindOneByEmail implements UserService.
func (u *UserService) FindOneByEmail(email string) (*entity.User, *domain.Error) {
	if email == "" {
		return nil, domain.NewError(400, errors.New("EMAIL_IS_EMPTY"), nil)
	}

	user, err := u.userRepo.FindOneByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindOneByID implements UserService.
func (u *UserService) FindOneByID(id int) (*entity.User, *domain.Error) {
	if id == 0 {
		return nil, domain.NewError(400, errors.New("ID_IS_EMPTY"), nil)
	}

	user, err := u.userRepo.FindOneByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// CreateOne implements UserService.
func (u *UserService) CreateOne(user entity.User) (*entity.User, *domain.Error) {
	newUser, err := u.userRepo.CreateOne(user)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

// UpdateOne implements UserService.
func (u *UserService) UpdateOneByID(user entity.User) (*entity.User, *domain.Error) {
	newUser, err := u.userRepo.UpdateOneByID(user.ID, user)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

// DeleteOneByID implements UserService.
func (u *UserService) DeleteOneByID(id int) *domain.Error {
	err := u.userRepo.DeleteOneByID(id)
	if err != nil {
		return err
	}
	return nil
}

func NewUserService(userRepo *UserRepo) *UserService {
	return &UserService{
		userRepo,
	}
}
