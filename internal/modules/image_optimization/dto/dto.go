package image_optimization_dto

type OptimizeImageRequest struct {
	Width   int    `json:"width" validate:"required"`
	Height  int    `json:"height" validate:"required"`
	Quality int    `json:"quality" validate:"omitempty"`
	Path    string `json:"path" validate:"required"`
}
