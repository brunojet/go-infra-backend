package dtos

type ApplicationVersionPost struct {
}

type ApplicationVersionPatch struct {
	Stage ApplicationStage `json:"stage" binding:"required,oneof=pilot production archived"`
}

// DTO para leitura de versão de aplicação
type ApplicationVersionGet struct {
	ApplicationVersionId       int64                `json:"applicationVersionId,string"`
	ApplicationConfigurationId string               `json:"applicationConfigurationId"`
	Stage                      ApplicationStage     `json:"stage"`
	ReviewAt                   string               `json:"reviewAt,omitempty"`
	ProductionAt               string               `json:"productionAt,omitempty"`
	ApplicationCatalogs        []ApplicationCatalog `json:"applicationCatalogs,omitempty"`
	BaseTimestamps             `json:",inline"`
}
