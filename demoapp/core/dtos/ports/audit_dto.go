package ports

import (
	"time"

	"github.com/brunojet/go-infra-backend/demoapp/core/models/ports"
)

type AuditDTO struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func ToAuditResponse(a ports.Audit) AuditDTO {
	return AuditDTO{
		CreatedAt: a.CreatedAt.UTC(),
		UpdatedAt: a.UpdatedAt.UTC(),
	}
}
