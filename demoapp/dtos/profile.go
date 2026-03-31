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
	ApplicationProfileId          int64                             `json:"applicationProfileId,string"`
	ApplicationId                 int64                             `json:"applicationId,string"`
	Name                          string                            `json:"name"`
	Description                   string                            `json:"description,omitempty"`
	Stage                         ApplicationStage                  `json:"stage"` //pending, review, production
	ReviewAt                      string                            `json:"reviewAt,omitempty"`
	ProductionAt                  string                            `json:"productionAt,omitempty"`
	Filters                       []FilterGet                       `json:"filters,omitempty"`
	Application                   *ApplicationGet                   `json:"application,omitempty"`
	Icon                          ApplicationImageGet               `json:"icon"`
	ApplicationCatalogs           []ApplicationCatalog              `json:"applicationCatalogs,omitempty"`
	ApplicationProfileScreenshots []ApplicationProfileScreenshotGet `json:"applicationProfileScreenshots,omitempty"`
	BaseTimestamps                `json:",inline"`
}

type ApplicationProfileScreenshotPostDTO struct {
	Position   int16                `json:"position" binding:"required"`
	Screenshot ApplicationImagePost `json:"screenshot" binding:"required"`
}

type ApplicationProfileScreenshotGet struct {
	Position       int16               `json:"position"`
	Screenshot     ApplicationImageGet `json:"screenshot"`
	BaseTimestamps `json:",inline"`
}
