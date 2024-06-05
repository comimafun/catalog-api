package user

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"

	"gorm.io/gorm"
)

type UserRepo interface {
	CreateOne(user entity.User) (*entity.User, *domain.Error)
	DeleteOneByID(id int) *domain.Error
	FindOneByID(id int) (*entity.User, *domain.Error)
	FindOneByEmail(email string) (*entity.User, *domain.Error)
	UpdateOneByID(id int, user entity.User) (*entity.User, *domain.Error)
}

type userRepo struct {
	db *gorm.DB
}

// DeleteOneByID implements UserRepo.
func (u *userRepo) DeleteOneByID(id int) *domain.Error {
	err := u.db.Delete(&entity.User{}, id).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// FindOneByEmail implements UserRepo.
func (u *userRepo) FindOneByEmail(email string) (*entity.User, *domain.Error) {
	var user entity.User
	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &user, nil
}

// FindOneByID implements UserRepo.
func (u *userRepo) FindOneByID(id int) (*entity.User, *domain.Error) {
	var user entity.User
	if err := u.db.First(&user, id).Error; err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &user, nil
}

// UpdateOneByID implements UserRepo.
func (u *userRepo) UpdateOneByID(id int, user entity.User) (*entity.User, *domain.Error) {
	var updated entity.User
	if err := u.db.Model(&entity.User{}).Where("id = ?", id).Updates(&user).Scan(&updated).Error; err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &updated, nil
}

// CreateOne implements UserRepo.
func (u *userRepo) CreateOne(user entity.User) (*entity.User, *domain.Error) {
	var newUser entity.User
	if err := u.db.Create(&user).Scan(&newUser).Error; err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &newUser, nil
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{
		db,
	}
}
