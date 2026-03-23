package dtos

type FilterTypeDTO struct {
	FilterTypeId int64  `json:"filterTypeId,string"`
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	CreatedAt    string `json:"createdAt,omitempty"`
	UpdatedAt    string `json:"updatedAt,omitempty"`
	DeletedAt    string `json:"deletedAt,omitempty"`
}

type FilterDTO struct {
	FilterId     int64         `json:"filterId,string"`
	FilterTypeId int64         `json:"filterTypeId,string"`
	Name         string        `json:"name"`
	Description  string        `json:"description,omitempty"`
	CreatedAt    string        `json:"createdAt,omitempty"`
	UpdatedAt    string        `json:"updatedAt,omitempty"`
	DeletedAt    string        `json:"deletedAt,omitempty"`
	FilterType   FilterTypeDTO `json:"filterType,omitempty"`
}
