package dtos

type FilterTypePostDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}

type FilterTypePatchDTO struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type FilterTypeGetDTO struct {
	FilterTypeId      int64  `json:"filterTypeId,string"`
	Name              string `json:"name"`
	Description       string `json:"description,omitempty"`
	BaseTimestampsDTO `json:",inline"`
}

type FilterPostDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}

type FilterPatchDTO struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type FilterGetDTO struct {
	FilterId          int64             `json:"filterId,string"`
	FilterTypeId      int64             `json:"filterTypeId,string"`
	Name              string            `json:"name"`
	Description       string            `json:"description,omitempty"`
	FilterType        *FilterTypeGetDTO `json:"filterType,omitempty"`
	BaseTimestampsDTO `json:",inline"`
}
