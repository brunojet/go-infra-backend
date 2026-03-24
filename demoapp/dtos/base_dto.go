package dtos

type BaseTimestampsDTO struct {
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	DeletedAt string `json:"deletedAt,omitempty"`
}

type BaseFileDTO struct {
	FileName    string `json:"fileName" binding:"required"`
	FileHash    string `json:"fileHash" binding:"required"`
	FileSize    int64  `json:"fileSize" binding:"required"`
	ContentType string `json:"contentType" binding:"required"`
}
