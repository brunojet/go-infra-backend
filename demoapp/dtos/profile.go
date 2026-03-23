package dtos

type ApplicationProfileDTO struct {
	ApplicationProfileId          int64                             `json:"applicationProfileId,string"`
	ApplicationId                 int64                             `json:"applicationId,string"`
	Stage                         int16                             `json:"stage"`
	Name                          string                            `json:"name"`
	Description                   string                            `json:"description,omitempty"`
	ApplicationImageId            int64                             `json:"applicationImageId,string"`
	ReviewAt                      string                            `json:"reviewAt,omitempty"`
	ProductionAt                  string                            `json:"productionAt,omitempty"`
	CreatedAt                     string                            `json:"createdAt,omitempty"`
	UpdatedAt                     string                            `json:"updatedAt,omitempty"`
	DeletedAt                     string                            `json:"deletedAt,omitempty"`
	Filters                       []FilterDTO                       `json:"filters,omitempty"`
	Application                   *ApplicationDTO                   `json:"application,omitempty"`
	ApplicationImage              *ApplicationImageDTO              `json:"applicationImage,omitempty"`
	ApplicationCatalogs           []ApplicationCatalogDTO           `json:"applicationCatalogs,omitempty"`
	ApplicationProfileScreenshots []ApplicationProfileScreenshotDTO `json:"applicationProfileScreenshots,omitempty"`
}

type ApplicationProfileScreenshotDTO struct {
	ApplicationProfileId int64                `json:"applicationProfileId,string"`
	ApplicationImageId   int64                `json:"applicationImageId,string"`
	Position             int16                `json:"position"`
	CreatedAt            string               `json:"createdAt,omitempty"`
	UpdatedAt            string               `json:"updatedAt,omitempty"`
	DeletedAt            string               `json:"deletedAt,omitempty"`
	ApplicationImage     *ApplicationImageDTO `json:"applicationImage,omitempty"`
}
