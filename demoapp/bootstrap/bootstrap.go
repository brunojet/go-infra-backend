package demoapp

import (
	hnd "github.com/brunojet/go-infra-backend/demoapp/handlers"
	repo "github.com/brunojet/go-infra-backend/demoapp/repositories"
	svc "github.com/brunojet/go-infra-backend/demoapp/services"
	"github.com/brunojet/go-infra-backend/pkg/database"
	"github.com/gin-gonic/gin"
)

func SetupHelloWorldModule(db database.DatabaseAdapter, rg *gin.RouterGroup) error {
	helloWorldRepo := repo.NewHelloWorldRepo(db)
	helloWorldService := svc.NewHelloWorldService(helloWorldRepo)
	hnd.NewHelloWorldHandler(rg, helloWorldService)
	return nil
}

func SetupStoreProviderModule(db database.DatabaseAdapter, rg *gin.RouterGroup) error {
	// Repositories
	filterTypeRepo := repo.NewFilterTypeRepo(db)
	filterRepo := repo.NewFilterRepo(db)
	terminalModelRepo := repo.NewTerminalModelRepo(db)
	terminalModelConfigurationRepo := repo.NewTerminalModelConfigurationRepo(db)
	applicationRepo := repo.NewApplicationRepo(db)
	applicationConfigurationRepo := repo.NewApplicationConfigurationRepo(db)
	applicationProfileRepo := repo.NewApplicationProfileRepo(db)
	applicationVersionRepo := repo.NewApplicationVersionRepo(db)
	applicationCatalogRepo := repo.NewApplicationCatalogRepo(db)

	// Services
	filterTypeService := svc.NewFilterTypeService(filterTypeRepo)
	filterNestedService := svc.NewFilterNestedService(filterRepo)
	terminalModelService := svc.NewTerminalModelService(terminalModelRepo)
	terminalModelConfigurationNestedService := svc.NewTerminalModelConfigurationNestedService(terminalModelConfigurationRepo)
	applicationService := svc.NewApplicationService(applicationRepo)
	applicationConfigurationService := svc.NewApplicationConfigurationNestedService(applicationConfigurationRepo)
	applicationProfileNestedService := svc.NewApplicationProfileNestedService(applicationProfileRepo, applicationConfigurationRepo, applicationVersionRepo, applicationCatalogRepo)
	applicationVersionNestedService := svc.NewApplicationVersionNestedService(applicationVersionRepo, applicationProfileRepo, applicationCatalogRepo)

	// Handlers
	hnd.NewFilterTypeHandler(rg, filterTypeService)
	hnd.NewFiltersHandler(rg, filterNestedService)
	hnd.NewTerminalModelHandler(rg, terminalModelService)
	hnd.NewTerminalModelConfigurationHandler(rg, terminalModelConfigurationNestedService)
	hnd.NewApplicationHandler(rg, applicationService)
	hnd.NewApplicationConfigurationHandler(rg, applicationConfigurationService)
	hnd.NewApplicationProfileHandler(rg, applicationProfileNestedService)
	hnd.NewApplicationVersionHandler(rg, applicationVersionNestedService)
	return nil
}
