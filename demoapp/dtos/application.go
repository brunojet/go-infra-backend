package dtos

type ApplicationPostDTO struct {
	CustomerId  string `json:"customerId" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}

type ApplicationPatchDTO struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type ApplicationGetDTO struct {
	ApplicationId             int64                            `json:"applicationId,string"`
	CustomerId                string                           `json:"customerId"`
	Name                      string                           `json:"name"`
	Description               string                           `json:"description,omitempty"`
	ApplicationConfigurations []ApplicationConfigurationGetDTO `json:"applicationConfigurations,omitempty"`
	BaseTimestampsDTO
}

type ApplicationConfigurationPostDTO struct {
	//Requerido junto ao parent applicationId para formar a PK composta.
	TerminalModelConfigurationId int64  `json:"terminalModelConfigurationId,string" binding:"required"`
	PackageName                  string `json:"packageName" binding:"required"`
}

type ApplicationConfigurationPatchDTO struct {
	PackageName string `json:"packageName,omitempty"`
}

type ApplicationConfigurationGetDTO struct {
	ApplicationConfigurationId string                            `json:"applicationConfigurationId"`
	PackageName                string                            `json:"packageName"`
	Application                *ApplicationGetDTO                `json:"application,omitempty"`
	TerminalModelConfiguration *TerminalModelConfigurationGetDTO `json:"terminalModelConfiguration,omitempty"`
	BaseTimestampsDTO          `json:",inline"`
}

// ApplicationCatalogDTO representa o catálogo de aplicações, com relacionamentos aninhados
// conforme o modelo ApplicationCatalog.
// Chave composta encodada em base64

type ApplicationCatalogGetDTO struct {
	ApplicationCatalogId     string                          `json:"applicationCatalogId"`
	ApplicationProfileId     int64                           `json:"applicationProfileId,string"`
	ApplicationVersionId     *int64                          `json:"applicationVersionId,string,omitempty"`
	ApplicationConfiguration *ApplicationConfigurationGetDTO `json:"applicationConfiguration,omitempty"`
	ApplicationProfile       *ApplicationProfileGetDTO       `json:"applicationProfile,omitempty"`
	ApplicationVersion       *ApplicationVersionGetDTO       `json:"applicationVersion,omitempty"`
	BaseTimestampsDTO        `json:",inline"`
}

// ApplicationImageDTO representa a imagem da aplicação, com relacionamentos aninhados
// conforme o modelo ApplicationImage.
// DTO para criação de imagem de aplicação (oneshot)
type ApplicationImagePostDTO struct {
	BaseFileDTO `json:",inline"`
}

type UploadPendingDTO struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
}

type DownloadReadyDTO struct {
	URL string `json:"url"`
}

type ApplicationImageGetDTO struct {
	ApplicationImageId int64 `json:"applicationImageId,string"`
	ApplicationId      int64 `json:"applicationId,string"`
	BaseFileDTO        `json:",inline"`
	State              string            `json:"state"`                   //upload_pending, processing, download_ready, failed
	UploadPendingDTO   *UploadPendingDTO `json:"uploadPending,omitempty"` //upload_pending, failed
	DownloadReadyDTO   *DownloadReadyDTO `json:"downloadReady,omitempty"` //download_ready
	BaseTimestampsDTO  `json:",inline"`
}
