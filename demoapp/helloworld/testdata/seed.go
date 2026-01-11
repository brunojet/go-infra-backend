package testdata

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/brunojet/go-infra-backend/demoapp/core/models"
	"gorm.io/gorm"
)

type HelloWorldSeed struct {
	Mensagem string `json:"mensagem"`
	Lingua   string `json:"lingua"`
}

func SeedHelloWorld(db *gorm.DB, basePath string) error {
	seedPath := filepath.Join(basePath, "helloworld", "testdata", "helloworld_seed.json")
	file, err := os.Open(seedPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var seeds []HelloWorldSeed
	if err := json.NewDecoder(file).Decode(&seeds); err != nil {
		return err
	}

	for _, s := range seeds {
		hw := models.HelloWorld{Message: s.Mensagem, Language: s.Lingua}
		db.Create(&hw)
	}
	return nil
}
