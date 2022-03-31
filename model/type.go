package model

type (
	// Metadata consists data meta responses
	Metadata struct {
		Page  int               `json:"page,omitempty"`
		Size  int               `json:"size,omitempty"`
		Order string            `json:"order,omitempty"`
		TotalData int           `json:"total_data,omitempty"`
		TotalPage int           `json:"total_page,omitempty"`
		Links map[string]string `json:"links,omitempty"`
	}

)