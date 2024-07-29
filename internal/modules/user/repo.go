package user

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

// DeleteOneByID implements UserRepo.
func (u *UserRepo) DeleteOneByID(id int) *domain.Error {
	err := u.db.Delete(&entity.User{}, id).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// FindOneByEmail implements UserRepo.
func (u *UserRepo) FindOneByEmail(email string) (*entity.User, *domain.Error) {
	var user entity.User
	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &user, nil
}

// FindOneByID implements UserRepo.
func (u *UserRepo) FindOneByID(id int) (*entity.User, *domain.Error) {
	var user entity.User
	if err := u.db.First(&user, id).Error; err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &user, nil
}

// UpdateOneByID implements UserRepo.
func (u *UserRepo) UpdateOneByID(id int, user entity.User) (*entity.User, *domain.Error) {
	var updated entity.User
	if err := u.db.Model(&entity.User{}).Where("id = ?", id).Updates(&user).Scan(&updated).Error; err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &updated, nil
}

// CreateOne implements UserRepo.
func (u *UserRepo) CreateOne(user entity.User) (*entity.User, *domain.Error) {
	if err := u.db.Create(&user).Error; err != nil {
		return nil, domain.NewError(500, err, nil)
	}
	return &user, nil
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db,
	}
}
