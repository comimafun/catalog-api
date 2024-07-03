package product_dto

type CreateProductBody struct {
	Name     string `json:"name" validate:"required,min=3,max=255"`
	ImageURL string `json:"image_url" validate:"required,url"`
}

type UpdateProductBody struct {
	ID int `json:"id" validate:"required,min=1"`
	CreateProductBody
}
