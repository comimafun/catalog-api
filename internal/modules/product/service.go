package product

import (
	"catalog-be/internal/domain"
	"catalog-be/internal/entity"
	product_dto "catalog-be/internal/modules/product/dto"
)

type ProductService interface {
	UpsertProductByCircleID(circleID int, inputs []product_dto.UpdateProductBody) ([]entity.Product, *domain.Error)
	UpsertOneProductByCircleID(circleID int, input product_dto.UpdateProductBody) (*entity.Product, *domain.Error)
	GetAllProductsByCircleID(circleID int) ([]entity.Product, *domain.Error)
	CountProductsByCircleID(circleID int) (int, *domain.Error)
}

type productService struct {
	repo ProductRepo
}

// CountProductsByCircleID implements ProductService.
func (p *productService) CountProductsByCircleID(circleID int) (int, *domain.Error) {
	return p.repo.CountProductsByCircleID(circleID)
}

// UpsertOneProductByCircleID implements ProductService.
func (p *productService) UpsertOneProductByCircleID(circleID int, input product_dto.UpdateProductBody) (*entity.Product, *domain.Error) {
	product := entity.Product{
		ID:       input.ID,
		CircleID: circleID,
		Name:     input.Name,
		ImageURL: input.ImageURL,
	}

	return p.repo.UpsertOneByCircleID(circleID, product)
}

// GetAllProductsByCircleID implements ProductService.
func (p *productService) GetAllProductsByCircleID(circleID int) ([]entity.Product, *domain.Error) {
	return p.repo.FindAllByCircleID(circleID)
}

// UpsertProductByCircleID implements ProductService.
func (p *productService) UpsertProductByCircleID(circleID int, inputs []product_dto.UpdateProductBody) ([]entity.Product, *domain.Error) {
	if len(inputs) == 0 {
		err := p.repo.DeleteAllByCircleID(circleID)
		if err != nil {
			return nil, err
		}
		return []entity.Product{}, nil
	}

	var products []entity.Product

	for _, input := range inputs {
		products = append(products, entity.Product{
			ID:       input.ID,
			CircleID: circleID,
			Name:     input.Name,
			ImageURL: input.ImageURL,
		})
	}

	return p.repo.BatchUpsertByCircleID(circleID, products)
}

func NewProductService(repo ProductRepo) ProductService {
	return &productService{
		repo: repo,
	}
}
