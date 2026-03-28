package dtos

type FilterTypePost struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}

type FilterTypePatch struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type FilterTypeGet struct {
	FilterTypeId   int64  `json:"filterTypeId,string"`
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	BaseTimestamps `json:",inline"`
}

type FilterPost struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}

type FilterPatch struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type FilterGet struct {
	FilterId       int64          `json:"filterId,string"`
	FilterTypeId   int64          `json:"filterTypeId,string"`
	Name           string         `json:"name"`
	Description    string         `json:"description,omitempty"`
	FilterType     *FilterTypeGet `json:"filterType,omitempty"`
	BaseTimestamps `json:",inline"`
}
