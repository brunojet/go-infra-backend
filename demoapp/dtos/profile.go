package dtos

type ApplicationProfilePost struct {
	Name                          string                                `json:"name" binding:"required"`
	Description                   string                                `json:"description,omitempty"`
	FilterIds                     []int64                               `json:"filterIds" binding:"required"`
	ApplicationImage              ApplicationImagePost                  `json:"applicationImage" binding:"required"`
	ApplicationProfileScreenshots []ApplicationProfileScreenshotPostDTO `json:"applicationProfileScreenshots,omitempty"`
}

type ApplicationProfilePatch struct {
	Stage ApplicationStage `json:"stage" binding:"required,oneof=review production archived"` //pending, review, production
}

type ApplicationProfileGet struct {
	ApplicationProfileId          int64                          `json:"applicationProfileId,string"`
	ApplicationId                 int64                          `json:"applicationId,string"`
	Name                          string                         `json:"name"`
	Description                   string                         `json:"description,omitempty"`
	Stage                         ApplicationStage               `json:"stage"` //pending, review, production
	ReviewAt                      string                         `json:"reviewAt,omitempty"`
	ProductionAt                  string                         `json:"productionAt,omitempty"`
	Filters                       []FilterGet                    `json:"filters,omitempty"`
	Application                   *ApplicationGet                `json:"application,omitempty"`
	ApplicationImage              ApplicationImage               `json:"applicationImage"`
	ApplicationCatalogs           []ApplicationCatalog           `json:"applicationCatalogs,omitempty"`
	ApplicationProfileScreenshots []ApplicationProfileScreenshot `json:"applicationProfileScreenshots,omitempty"`
	BaseTimestamps                `json:",inline"`
}

type ApplicationProfileScreenshotPostDTO struct {
	Position             int16 `json:"position" binding:"required"`
	ApplicationImagePost `json:",inline" binding:"required"`
}

type ApplicationProfileScreenshot struct {
	ApplicationProfileId int64            `json:"applicationProfileId,string"`
	ApplicationImageId   int64            `json:"applicationImageId,string"`
	Position             int16            `json:"position"`
	ApplicationImage     ApplicationImage `json:"applicationImage"`
	BaseTimestamps       `json:",inline"`
}
