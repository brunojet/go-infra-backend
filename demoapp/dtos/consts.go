package dtos

type FileStatus string

const (
	ApplicationImageStatusPending    FileStatus = "pending"
	ApplicationImageStatusProcessing FileStatus = "processing"
	ApplicationImageStatusReady      FileStatus = "ready"
	ApplicationImageStatusFailed     FileStatus = "failed"
)

type ApplicationStage string

const (
	ApplicationStagePending    ApplicationStage = "pending"
	ApplicationStageReview     ApplicationStage = "review"
	ApplicationStagePilot      ApplicationStage = "pilot"
	ApplicationStageProduction ApplicationStage = "production"
	ApplicationStageArchived   ApplicationStage = "archived"
)

type IntegrationType string

const (
	IntegrationTypeRFAL IntegrationType = "RFAL"
	IntegrationTypeTEF  IntegrationType = "TEF"
)
