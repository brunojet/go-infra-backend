package repositories

import "github.com/brunojet/go-infra-backend/demoapp/models"

const (
	whereApplicationIDEq   = models.ColApplicationID + " = ?"
	whereTerminalModelIDEq = models.ColTerminalModelConfigurationID + " = ?"
	whereStageEq           = models.ColStage + " = ?"
)
