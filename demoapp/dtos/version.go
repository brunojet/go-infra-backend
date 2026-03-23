package dtos

// DTO para criação de versão de aplicação
type ApplicationVersionPostDTO struct {
	ApplicationConfigurationId string `json:"applicationConfigurationId" binding:"required"`
}

// DTO para atualização de versão de aplicação
type ApplicationVersionPatchDTO struct {
	Stage string `json:"stage" binding:"required,oneof=pilot production archived"` //pilot, production, archived
}

// DTO para leitura de versão de aplicação
type ApplicationVersionGetDTO struct {
	ApplicationVersionId       int64                      `json:"applicationVersionId,string"`
	ApplicationConfigurationId string                     `json:"applicationConfigurationId"`
	Stage                      string                     `json:"stage"` //pending, pilot, production, archived
	ReviewAt                   string                     `json:"reviewAt,omitempty"`
	ProductionAt               string                     `json:"productionAt,omitempty"`
	ApplicationCatalogs        []ApplicationCatalogGetDTO `json:"applicationCatalogs,omitempty"`
	BaseTimestampsDTO          `json:",inline"`
}
