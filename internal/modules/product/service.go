package product

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	"errors"

	"gorm.io/gorm"
)

type ProductService struct {
	repo *ProductRepo
}

// DeleteOneProductByID implements ProductService.
func (p *ProductService) DeleteOneProductByID(circleID int, id int) *domain.Error {
	return p.repo.DeleteOneProductByProductID(circleID, id)
}

// UpdateOneProductByCircleAndProductID implements ProductService.
func (p *ProductService) UpdateOneProductByCircleAndProductID(circleID int, input entity.Product) (*entity.Product, *domain.Error) {
	check, err := p.repo.GetOneProductByProductID(input.ID)
	if err != nil {
		if errors.Is(err.Err, gorm.ErrRecordNotFound) {
			return nil, domain.NewError(404, errors.New("PRODUCT_NOT_FOUND"), nil)
		}

		return nil, err
	}

	if check.CircleID != circleID {
		return nil, domain.NewError(403, errors.New("FORBIDDEN"), nil)
	}

	return p.repo.UpdateOneByProductID(input.ID, input)
}

// CountProductsByCircleID implements ProductService.
func (p *ProductService) CountProductsByCircleID(circleID int) (int, *domain.Error) {
	return p.repo.CountProductsByCircleID(circleID)
}

// CreateOneProductByCircleID implements ProductService.
func (p *ProductService) CreateOneProductByCircleID(circleID int, input entity.Product) (*entity.Product, *domain.Error) {
	product := entity.Product{
		ID:       input.ID,
		CircleID: circleID,
		Name:     input.Name,
		ImageURL: input.ImageURL,
	}

	return p.repo.CreateOneOneByCircleID(circleID, product)
}

// GetAllProductsByCircleID implements ProductService.
func (p *ProductService) GetAllProductsByCircleID(circleID int) ([]entity.Product, *domain.Error) {
	return p.repo.GetAllProductByCircleID(circleID)
}

func NewProductService(repo *ProductRepo) *ProductService {
	return &ProductService{
		repo: repo,
	}
}
