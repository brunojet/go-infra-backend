package dtos

type TerminalModelPostDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}

type TerminalModelPatchDTO struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type TerminalModelGetDTO struct {
	TerminalModelId             int64                              `json:"terminalModelId,string"`
	Name                        string                             `json:"name"`
	Description                 string                             `json:"description,omitempty"`
	TerminalModelConfigurations []TerminalModelConfigurationGetDTO `json:"terminalModelConfigurations,omitempty"`
	BaseTimestampsDTO           `json:",inline"`
}

type TerminalModelConfigurationPostDTO struct {
	IntegrationType int16 `json:"integrationType" binding:"required"`
}

type TerminalModelConfigurationPatchDTO struct {
	IntegrationType int16 `json:"integrationType,omitempty"`
}

type TerminalModelConfigurationGetDTO struct {
	TerminalModelConfigurationId int64                `json:"terminalModelConfigurationId,string"`
	IntegrationType              int16                `json:"integrationType"`
	TerminalModel                *TerminalModelGetDTO `json:"terminalModel,omitempty"`
	BaseTimestampsDTO            `json:",inline"`
}
