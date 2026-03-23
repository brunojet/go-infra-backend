package dtos

type ApplicationVersionDTO struct {
	ApplicationVersionId int64 `json:"applicationVersionId,string"`
	// Chave composta de ApplicationConfiguration encodada em base64
	ApplicationConfigurationId string                  `json:"applicationConfigurationId"`
	Stage                      int16                   `json:"stage"`
	ReviewAt                   string                  `json:"reviewAt,omitempty"`
	ProductionAt               string                  `json:"productionAt,omitempty"`
	CreatedAt                  string                  `json:"createdAt,omitempty"`
	UpdatedAt                  string                  `json:"updatedAt,omitempty"`
	DeletedAt                  string                  `json:"deletedAt,omitempty"`
	ApplicationCatalogs        []ApplicationCatalogDTO `json:"applicationCatalogs,omitempty"`
}
