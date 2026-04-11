package demobff

import (
	hnd "github.com/brunojet/go-infra-backend/demobff/handlers"
	svc "github.com/brunojet/go-infra-backend/demobff/services"
	"github.com/brunojet/go-infra-backend/pkg/infra/bffclient"
	bffrepo "github.com/brunojet/go-infra-backend/pkg/ports/bff/repositories"
	bffsvc "github.com/brunojet/go-infra-backend/pkg/ports/bff/services"
	"github.com/gin-gonic/gin"
)

func SetupHelloWorldBffModule(client bffclient.BffClient, rg *gin.RouterGroup) error {
	repo := bffrepo.NewBffRepository[svc.HelloWorldBffE, svc.HelloWorldBffE, svc.HelloWorldBffE](client)
	service := bffsvc.NewBffServiceImpl(repo, svc.HelloWorldBffMapper{})
	hnd.NewHelloWorldBffHandler(rg, service)
	return nil
}
