package dto

type Metadata struct {
	Page        int  `json:"page"`
	Limit       int  `json:"limit"`
	TotalDocs   int  `json:"total_docs"`
	TotalPages  int  `json:"total_pages"`
	HasNextPage bool `json:"has_next_page"`
}

type Pagination[T any] struct {
	Data     T        `json:"data"`
	Metadata Metadata `json:"metadata"`
}
