package product

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"

	"gorm.io/gorm"
)

type ProductRepo interface {
	FindAllByCircleID(circleID int) ([]entity.Product, *domain.Error)
	DeleteOneByID(id int) *domain.Error
	CreateOne(product entity.Product) (*entity.Product, *domain.Error)
	UpdateOneByID(id int, product entity.Product) (*entity.Product, *domain.Error)
	BatchUpsertByCircleID(circleID int, inputs []entity.Product) ([]entity.Product, *domain.Error)
}

type productRepo struct {
	db *gorm.DB
}

// BatchUpsertByCircleID implements ProductRepo.
func (p *productRepo) BatchUpsertByCircleID(circleID int, inputs []entity.Product) ([]entity.Product, *domain.Error) {
	tx := p.db.Begin()
	if tx.Error != nil {
		return nil, domain.NewError(500, tx.Error, nil)
	}
	var existingProducts []entity.Product
	err := tx.Where("circle_id = ?", circleID).Find(&existingProducts).Error
	if err != nil {
		tx.Rollback()
		return nil, domain.NewError(500, err, nil)
	}

	newProductMap := make(map[int]bool)

	for _, input := range inputs {
		if input.ID == 0 {
			err := tx.Create(&input).Error
			if err != nil {
				tx.Rollback()
				return nil, domain.NewError(500, err, nil)
			}

			newProductMap[input.ID] = true
		} else {
			err := tx.Model(&entity.Product{}).Where("id = ?", input.ID).Updates(input).Error
			if err != nil {
				tx.Rollback()
				return nil, domain.NewError(500, err, nil)
			}
			newProductMap[input.ID] = true
		}
	}

	for _, product := range existingProducts {
		_, ok := newProductMap[product.ID]
		if !ok {
			err := tx.Delete(&entity.Product{}, product.ID).Error
			if err != nil {
				tx.Rollback()
				return nil, domain.NewError(500, err, nil)
			}
		}
	}

	var updatedProducts []entity.Product

	err = tx.Where("circle_id = ?", circleID).Find(&updatedProducts).Error

	if err != nil {
		tx.Rollback()
		return nil, domain.NewError(500, err, nil)
	}

	tx.Commit()

	return updatedProducts, nil
}

// CreateOne implements ProductRepo.
func (p *productRepo) CreateOne(product entity.Product) (*entity.Product, *domain.Error) {
	err := p.db.Create(&product).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return &product, nil
}

// DeleteOneByID implements ProductRepo.
func (p *productRepo) DeleteOneByID(id int) *domain.Error {
	err := p.db.Delete(&entity.Product{}, id).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}

	return nil
}

// FindAllByCircleID implements ProductRepo.
func (p *productRepo) FindAllByCircleID(circleID int) ([]entity.Product, *domain.Error) {
	var products []entity.Product
	err := p.db.Where("circle_id = ?", circleID).Find(&products).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return products, nil
}

// UpdateOneByID implements ProductRepo.
func (p *productRepo) UpdateOneByID(id int, product entity.Product) (*entity.Product, *domain.Error) {
	err := p.db.Model(&entity.Product{}).Where("id = ?", id).Updates(product).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return &product, nil
}

func NewProductRepo(
	db *gorm.DB,
) ProductRepo {
	return &productRepo{
		db: db,
	}
}
