package upload_dto

type UploadImageBody struct {
	Type string `json:"type" validate:"required, oneof=products users circles"`
}
