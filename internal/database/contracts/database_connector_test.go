package contracts

import (
	"testing"

	"gorm.io/gorm"
)

type dummyConnector struct{}

func (dummyConnector) Open() (*gorm.DB, error) { return nil, nil }
func (dummyConnector) Close() error            { return nil }

func TestDatabaseConnectorInterface(t *testing.T) {
	var _ DatabaseConnector = dummyConnector{}
}
