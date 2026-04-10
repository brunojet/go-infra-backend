package dtos

type BaseTimestamps struct {
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	DeletedAt string `json:"deletedAt,omitempty"`
}

type BaseFilePost struct {
	FileName    string `json:"fileName" binding:"required"`
	FileHash    string `json:"fileHash" binding:"required"`
	FileSize    int64  `json:"fileSize" binding:"required"`
	ContentType string `json:"contentType" binding:"required"`
}

type BaseFile struct {
	BaseFilePost `json:",inline"`
	FileStatus   FileStatus `json:"fileStatus"`
}
