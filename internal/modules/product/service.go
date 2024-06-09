package product

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
)

type ProductService interface {
	UpsertProductByCircleID(circleID int, inputs []entity.Product) ([]entity.Product, *domain.Error)
}

type productService struct {
	repo ProductRepo
}

// UpsertProductByCircleID implements ProductService.
func (p *productService) UpsertProductByCircleID(circleID int, inputs []entity.Product) ([]entity.Product, *domain.Error) {
	return p.repo.BatchUpsertByCircleID(circleID, inputs)
}

func NewProductService(repo ProductRepo) ProductService {
	return &productService{
		repo: repo,
	}
}
