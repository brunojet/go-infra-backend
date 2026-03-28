package dtos

type TerminalModelPost struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description,omitempty"`
}

type TerminalModelPatch struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type TerminalModelGet struct {
	TerminalModelId             int64                           `json:"terminalModelId,string"`
	Name                        string                          `json:"name"`
	Description                 string                          `json:"description,omitempty"`
	TerminalModelConfigurations []TerminalModelConfigurationGet `json:"terminalModelConfigurations,omitempty"`
	BaseTimestamps              `json:",inline"`
}

type TerminalModelConfigurationPost struct {
	IntegrationType int16 `json:"integrationType" binding:"required"`
}

type TerminalModelConfigurationPatch struct {
	IntegrationType int16 `json:"integrationType,omitempty"`
}

type TerminalModelConfigurationGet struct {
	TerminalModelConfigurationId int64             `json:"terminalModelConfigurationId,string"`
	TerminalModelId              int64             `json:"terminalModelId,string"`
	IntegrationType              int16             `json:"integrationType"`
	TerminalModel                *TerminalModelGet `json:"terminalModel,omitempty"`
	BaseTimestamps               `json:",inline"`
}
