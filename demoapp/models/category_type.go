package models

import (
	"database/sql"
)

// CategoryType minimal definition to satisfy references from Category.
type CategoryType struct {
	CategoryTypeId int64          `gorm:"primaryKey;autoIncrement" json:"category_type_id"`
	Name           sql.NullString `gorm:"size:128;index:ux_category_type_name" json:"name,omitempty"`
}

func (CategoryType) TableName() string { return "category_type" }
