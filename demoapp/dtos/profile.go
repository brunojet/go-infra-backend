package dtos

type ApplicationProfilePostDTO struct {
	Name                          string                                `json:"name" binding:"required"`
	Description                   string                                `json:"description,omitempty"`
	FilterIds                     []int64                               `json:"filterIds" binding:"required"`
	ApplicationImage              ApplicationImagePostDTO               `json:"applicationImage" binding:"required"`
	ApplicationProfileScreenshots []ApplicationProfileScreenshotPostDTO `json:"applicationProfileScreenshots,omitempty"`
}

type ApplicationProfilePatchDTO struct {
	Stage string `json:"stage" binding:"required,oneof=review production archived"` //pending, review, production
}

type ApplicationProfileGetDTO struct {
	ApplicationProfileId          int64                                `json:"applicationProfileId,string"`
	ApplicationId                 int64                                `json:"applicationId,string"`
	Name                          string                               `json:"name"`
	Description                   string                               `json:"description,omitempty"`
	Stage                         string                               `json:"stage"` //pending, review, production
	ReviewAt                      string                               `json:"reviewAt,omitempty"`
	ProductionAt                  string                               `json:"productionAt,omitempty"`
	Filters                       []FilterGetDTO                       `json:"filters,omitempty"`
	Application                   *ApplicationGetDTO                   `json:"application,omitempty"`
	ApplicationImage              ApplicationImageGetDTO               `json:"applicationImage"`
	ApplicationCatalogs           []ApplicationCatalogGetDTO           `json:"applicationCatalogs,omitempty"`
	ApplicationProfileScreenshots []ApplicationProfileScreenshotGetDTO `json:"applicationProfileScreenshots,omitempty"`
	BaseTimestampsDTO             `json:",inline"`
}

type ApplicationProfileScreenshotPostDTO struct {
	Position                int16 `json:"position" binding:"required"`
	ApplicationImagePostDTO `json:",inline" binding:"required"`
}

type ApplicationProfileScreenshotGetDTO struct {
	ApplicationProfileId int64                  `json:"applicationProfileId,string"`
	ApplicationImageId   int64                  `json:"applicationImageId,string"`
	Position             int16                  `json:"position"`
	ApplicationImage     ApplicationImageGetDTO `json:"applicationImage"`
	BaseTimestampsDTO    `json:",inline"`
}
