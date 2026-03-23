package dtos

type ApplicationDTO struct {
	ApplicationId             int64                         `json:"applicationId,string"`
	CustomerId                string                        `json:"customerId"`
	Name                      string                        `json:"name"`
	Description               string                        `json:"description,omitempty"`
	CreatedAt                 string                        `json:"createdAt,omitempty"`
	UpdatedAt                 string                        `json:"updatedAt,omitempty"`
	DeletedAt                 string                        `json:"deletedAt,omitempty"`
	ApplicationConfigurations []ApplicationConfigurationDTO `json:"applicationConfigurations,omitempty"`
}

type ApplicationConfigurationDTO struct {
	// Chave composta encodada em base64
	ApplicationConfigurationId   string                         `json:"applicationConfigurationId"`
	TerminalModelConfigurationId int64                          `json:"terminalModelConfigurationId,string"`
	PackageName                  string                         `json:"packageName"`
	CreatedAt                    string                         `json:"createdAt,omitempty"`
	UpdatedAt                    string                         `json:"updatedAt,omitempty"`
	DeletedAt                    string                         `json:"deletedAt,omitempty"`
	Application                  *ApplicationDTO                `json:"application,omitempty"`
	TerminalModelConfiguration   *TerminalModelConfigurationDTO `json:"terminalModelConfiguration,omitempty"`
}

// ApplicationCatalogDTO representa o catálogo de aplicações, com relacionamentos aninhados
// conforme o modelo ApplicationCatalog.
// Chave composta encodada em base64
type ApplicationCatalogDTO struct {
	ApplicationCatalogId string                 `json:"applicationCatalogId"`
	ApplicationProfileId int64                  `json:"applicationProfileId,string"`
	ApplicationVersionId *int64                 `json:"applicationVersionId,string,omitempty"`
	CreatedAt            string                 `json:"createdAt,omitempty"`
	UpdatedAt            string                 `json:"updatedAt,omitempty"`
	DeletedAt            string                 `json:"deletedAt,omitempty"`
	ApplicationProfile   *ApplicationProfileDTO `json:"applicationProfile,omitempty"`
	ApplicationVersion   *ApplicationVersionDTO `json:"applicationVersion,omitempty"`
}

// ApplicationImageDTO representa a imagem da aplicação, com relacionamentos aninhados
// conforme o modelo ApplicationImage.
type ApplicationImageDTO struct {
	ApplicationImageId int64  `json:"applicationImageId,string"`
	ApplicationId      int64  `json:"applicationId,string"`
	FileName           string `json:"fileName,omitempty"`
	FileContentType    string `json:"fileContentType,omitempty"`
	FileHash           []byte `json:"fileHash"`
	ImageType          *int16 `json:"imageType,omitempty"`
	CreatedAt          string `json:"createdAt,omitempty"`
	UpdatedAt          string `json:"updatedAt,omitempty"`
	DeletedAt          string `json:"deletedAt,omitempty"`
}
