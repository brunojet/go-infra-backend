package repositories

import "errors"

// Erros de domínio específicos dos repositórios de demoapp
var (
	ErrInvalidApplicationProfileID         = errors.New("application_profile_id must be valid")
	ErrInvalidApplicationID                = errors.New("application_id must be valid")
	ErrInvalidConfigurationTerminalModelID = errors.New("terminal_model_id must be valid")
	ErrInvalidApplicationVersionID         = errors.New("application_version_id must be valid")
	ErrInvalidStage                        = errors.New("stage is required")

	ErrInvalidApplicationVersionModel = errors.New("application version model is required")
	ErrInvalidApplicationProfileModel = errors.New("application profile model is required")

	ErrCustomerIdRequired     = errors.New("customer_id is required")
	ErrCustomerIdViolation    = errors.New("customer_id violates ownership rules")
	ErrWhereArgInvalid        = errors.New("where argument is invalid")
	ErrApplicationIdViolation = errors.New("application_id violates ownership rules")
)
