package product

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"

	"gorm.io/gorm"
)

type ProductRepo struct {
	db *gorm.DB
}

// GetOneProductByProductID implements ProductRepo.
func (p *ProductRepo) GetOneProductByProductID(id int) (*entity.Product, *domain.Error) {
	var product entity.Product
	err := p.db.First(&product, id).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return &product, nil
}

// CountProductsByCircleID implements ProductRepo.
func (p *ProductRepo) CountProductsByCircleID(circleID int) (int, *domain.Error) {
	var count int64
	err := p.db.Model(&entity.Product{}).Where("circle_id = ?", circleID).Count(&count).Error
	if err != nil {
		return 0, domain.NewError(500, err, nil)
	}

	return int(count), nil
}

// CreateOneOneByCircleID implements ProductRepo.
func (p *ProductRepo) CreateOneOneByCircleID(circleID int, product entity.Product) (*entity.Product, *domain.Error) {
	err := p.db.Where("circle_id = ?", circleID).Save(&product).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return &product, nil
}

// DeleteAllByCircleID implements ProductRepo.
func (p *ProductRepo) DeleteAllByCircleID(circleID int) *domain.Error {
	err := p.db.Where("circle_id = ?", circleID).Delete(&entity.Product{}).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}
	return nil
}

// BatchUpsertByCircleID implements ProductRepo.
func (p *ProductRepo) BatchUpsertByCircleID(circleID int, inputs []entity.Product) ([]entity.Product, *domain.Error) {
	tx := p.db.Begin()
	if tx.Error != nil {
		return nil, domain.NewError(500, tx.Error, nil)
	}

	createdOrUpdatedProductsIDs := make(map[int]bool)

	for _, input := range inputs {
		if input.ID == 0 {
			err := tx.Create(&input).Error
			if err != nil {
				tx.Rollback()
				return nil, domain.NewError(500, err, nil)
			}

			createdOrUpdatedProductsIDs[input.ID] = true
		} else {
			err := tx.Model(&entity.Product{}).Where("id = ? AND circle_id", input.ID, circleID).Updates(input).Error
			if err != nil {
				tx.Rollback()
				return nil, domain.NewError(500, err, nil)
			}
			createdOrUpdatedProductsIDs[input.ID] = true
		}
	}

	var previousProducts []entity.Product
	err := tx.Where("circle_id = ?", circleID).Find(&previousProducts).Error
	if err != nil {
		tx.Rollback()
		return nil, domain.NewError(500, err, nil)
	}

	var idsToDelete []int
	for _, product := range previousProducts {
		_, ok := createdOrUpdatedProductsIDs[product.ID]
		if !ok {
			idsToDelete = append(idsToDelete, product.ID)
		}
	}

	if len(idsToDelete) > 0 {
		err := tx.Where("id IN (?) AND circle_id = ?", idsToDelete, circleID).Delete(&entity.Product{}).Error
		if err != nil {
			tx.Rollback()
			return nil, domain.NewError(500, err, nil)
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

// CreateOneProduct implements ProductRepo.
func (p *ProductRepo) CreateOneProduct(product entity.Product) (*entity.Product, *domain.Error) {
	err := p.db.Create(&product).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return &product, nil
}

// DeleteOneProductByProductID implements ProductRepo.
func (p *ProductRepo) DeleteOneProductByProductID(circleID int, id int) *domain.Error {
	err := p.db.Delete(&entity.Product{}, id).Error
	if err != nil {
		return domain.NewError(500, err, nil)
	}

	return nil
}

// GetAllProductByCircleID implements ProductRepo.
func (p *ProductRepo) GetAllProductByCircleID(circleID int) ([]entity.Product, *domain.Error) {
	var products []entity.Product
	err := p.db.Where("circle_id = ?", circleID).Find(&products).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return products, nil
}

// UpdateOneByProductID implements ProductRepo.
func (p ProductRepo) UpdateOneByProductID(id int, product entity.Product) (*entity.Product, *domain.Error) {
	err := p.db.Model(&entity.Product{}).Where("id = ?", id).Updates(&product).Error
	if err != nil {
		return nil, domain.NewError(500, err, nil)
	}

	return &product, nil
}

func NewProductRepo(
	db *gorm.DB,
) *ProductRepo {
	return &ProductRepo{
		db: db,
	}
}
