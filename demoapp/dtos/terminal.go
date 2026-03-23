package dtos

type TerminalModelDTO struct {
	TerminalModelId             int64                           `json:"terminalModelId,string"`
	Name                        string                          `json:"name"`
	Description                 string                          `json:"description,omitempty"`
	CreatedAt                   string                          `json:"createdAt,omitempty"`
	UpdatedAt                   string                          `json:"updatedAt,omitempty"`
	DeletedAt                   string                          `json:"deletedAt,omitempty"`
	TerminalModelConfigurations []TerminalModelConfigurationDTO `json:"terminalModelConfigurations,omitempty"`
}

type TerminalModelConfigurationDTO struct {
	TerminalModelConfigurationId int64             `json:"terminalModelConfigurationId,string"`
	TerminalModelId              int64             `json:"terminalModelId,string"`
	IntegrationType              int16             `json:"integrationType"`
	CreatedAt                    string            `json:"createdAt,omitempty"`
	UpdatedAt                    string            `json:"updatedAt,omitempty"`
	DeletedAt                    string            `json:"deletedAt,omitempty"`
	TerminalModel                *TerminalModelDTO `json:"terminalModel,omitempty"`
}
