package contracts

import "gorm.io/gorm"

type DatabaseConnector interface {
	Open() (*gorm.DB, error)
	Close() error
}
