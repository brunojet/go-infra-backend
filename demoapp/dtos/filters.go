package dtos

type FilterTypeDTO struct {
	FilterTypeId int64  `json:"filterTypeId"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	CreatedAt    string `json:"createdAt,omitempty"`
	UpdatedAt    string `json:"updatedAt,omitempty"`
	DeletedAt    string `json:"deletedAt,omitempty"`
}

type FilterDTO struct {
	FilterId     int64         `json:"filterId"`
	FilterTypeId int64         `json:"filterTypeId"`
	Name         string        `json:"name"`
	Description  string        `json:"description,omitempty"`
	CreatedAt    string        `json:"createdAt,omitempty"`
	UpdatedAt    string        `json:"updatedAt,omitempty"`
	DeletedAt    string        `json:"deletedAt,omitempty"`
	FilterType   FilterTypeDTO `json:"filterType,omitempty"`
}
