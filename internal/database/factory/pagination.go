package factory

import (
	"catalog-be/internal/dto"
	"math"
)

func GetPaginationMetadata(totalDocs int, page int, limit int) *dto.Metadata {
	totalPages := int(math.Ceil(float64(totalDocs) / float64(limit)))
	return &dto.Metadata{
		TotalDocs:   totalDocs,
		TotalPages:  totalPages,
		HasNextPage: page < totalPages,
		Page:        page,
		Limit:       limit,
	}
}
