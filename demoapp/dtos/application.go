package dtos

type ApplicationPost struct {
	CustomerId  string `json:"customerId" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}

type ApplicationPatch struct {
	CustomerId  string `json:"customerId,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type ApplicationGet struct {
	ApplicationId             int64                         `json:"applicationId,string"`
	CustomerId                string                        `json:"customerId"`
	Name                      string                        `json:"name"`
	Description               string                        `json:"description,omitempty"`
	ApplicationConfigurations []ApplicationConfigurationGet `json:"applicationConfigurations,omitempty"`
	BaseTimestamps
}

type ApplicationConfigurationPost struct {
	//Requerido junto ao parent applicationId para formar a PK composta.
	TerminalModelConfigurationId int64  `json:"terminalModelConfigurationId,string" binding:"required"`
	PackageName                  string `json:"packageName" binding:"required"`
}

type ApplicationConfigurationPatch struct {
	PackageName string `json:"packageName,omitempty"`
}

type ApplicationConfigurationGet struct {
	ApplicationConfigurationId string                         `json:"applicationConfigurationId"`
	PackageName                string                         `json:"packageName"`
	Application                *ApplicationGet                `json:"application,omitempty"`
	TerminalModelConfiguration *TerminalModelConfigurationGet `json:"terminalModelConfiguration,omitempty"`
	BaseTimestamps             `json:",inline"`
}

// ApplicationCatalogDTO representa o catálogo de aplicações, com relacionamentos aninhados
// conforme o modelo ApplicationCatalog.
// Chave composta encodada em base64

type ApplicationCatalog struct {
	ApplicationCatalogId     string                       `json:"applicationCatalogId"`
	ApplicationProfileId     int64                        `json:"applicationProfileId,string"`
	ApplicationVersionId     *int64                       `json:"applicationVersionId,string,omitempty"`
	ApplicationConfiguration *ApplicationConfigurationGet `json:"applicationConfiguration,omitempty"`
	ApplicationProfile       *ApplicationProfileGet       `json:"applicationProfile,omitempty"`
	ApplicationVersion       *ApplicationVersionGet       `json:"applicationVersion,omitempty"`
	BaseTimestamps           `json:",inline"`
}

// ApplicationImageDTO representa a imagem da aplicação, com relacionamentos aninhados
// conforme o modelo ApplicationImage.
// DTO para criação de imagem de aplicação (oneshot)
type ApplicationImagePost struct {
	BaseFilePost `json:",inline"`
}

type UploadPending struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
}

type DownloadReady struct {
	URL string `json:"url"`
}

type ApplicationImageGet struct {
	UploadPendingDTO *UploadPending `json:"uploadPending,omitempty"` //upload_pending, failed
	DownloadReadyDTO *DownloadReady `json:"downloadReady,omitempty"` //download_ready
	BaseFile         `json:",inline"`
	BaseTimestamps   `json:",inline"`
}
