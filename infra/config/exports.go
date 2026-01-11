package config

import (
	"github.com/brunojet/go-infra-backend/internal/config/adapters/env"
	cfgcontracts "github.com/brunojet/go-infra-backend/internal/config/contracts"
	cfgports "github.com/brunojet/go-infra-backend/internal/config/ports"
)

var NewEnvSource = env.New

type ConfigSource = cfgcontracts.Source

var (
	Trimmed          = cfgports.Trimmed
	Duration         = cfgports.Duration
	Bool             = cfgports.Bool
	Int              = cfgports.Int
	SplitCSV         = cfgports.SplitCSV
	SplitCSVUpper    = cfgports.SplitCSVUpper
	ParseKeyValueCSV = cfgports.ParseKeyValueCSV
)
