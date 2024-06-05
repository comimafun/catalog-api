package fandom_dto

type FindAllFilter struct {
	Search string `json:"search"`
	Page   int    `json:"page" validate:"required,min=1"`
	Limit  int    `json:"limit" validate:"required,min=1,max=20"`
}

type CreateBody struct {
	Name    string `json:"name" validate:"required,min=1,max=255"`
	Visible bool   `json:"visible" validate:"boolean"`
}
