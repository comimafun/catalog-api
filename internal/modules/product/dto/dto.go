package product_dto

type CreateUpdateProductBody struct {
	Name     string `json:"name" validate:"required,min=3,max=255"`
	ImageURL string `json:"image_url" validate:"required,url"`
}
