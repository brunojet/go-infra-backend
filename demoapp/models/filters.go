package models

import (
	"database/sql"

	repositories "github.com/brunojet/go-infra-backend/pkg/ports/backend/repositories"
	"gorm.io/gorm"
)

const (
	tableFilterType = "filter_type"
	tableFilter     = "filter"

	colFilterTypeID = "filter_type_id"
	colFilterName   = "name"
)

type FilterType struct {
	FilterTypeId int64          `gorm:"primaryKey;autoIncrement"`
	Name         sql.NullString `gorm:"size:128;uniqueIndex:ux_filter_type_name"`
	Description  sql.NullString `gorm:"size:500"`
	CreatedAt    sql.NullTime   `gorm:"autoCreateTime;index:idx_filter_type_del_created,priority:2"`
	UpdatedAt    sql.NullTime   `gorm:"autoUpdateTime;index:idx_filter_type_del_updated,priority:2"`
	DeletedAt    gorm.DeletedAt `gorm:"index:idx_filter_type_del_created,priority:1;index:idx_filter_type_del_updated,priority:1"`
	Filters      []Filter       `gorm:"foreignKey:FilterTypeId;references:FilterTypeId;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`
}

func (FilterType) TableName() string { return tableFilterType }

func (f FilterType) BeforeCreate(tx *gorm.DB) error {
	return repositories.AddOnConflictDoNothing(tx, colFilterName)
}

func (f FilterType) WhereOnConflict(tx *gorm.DB) *gorm.DB {
	return repositories.WhereOnConflict(tx,
		repositories.ConflictScope{ColumnName: colFilterName, ColumnValue: f.Name.String})
}

type Filter struct {
	FilterId     int64          `gorm:"primaryKey;autoIncrement"`
	FilterTypeId int64          `gorm:"not null;index:idx_filter_type;uniqueIndex:ux_filter_name_type,priority:2"`
	Name         sql.NullString `gorm:"size:128;uniqueIndex:ux_filter_name_type,priority:1"`
	Description  sql.NullString `gorm:"size:500"`
	CreatedAt    sql.NullTime   `gorm:"autoCreateTime;index:idx_filter_del_created,priority:2"`
	UpdatedAt    sql.NullTime   `gorm:"autoUpdateTime;index:idx_filter_del_updated,priority:2"`
	DeletedAt    gorm.DeletedAt `gorm:"index:idx_filter_del_created,priority:1;index:idx_filter_del_updated,priority:1"`

	FilterType *FilterType
}

func (Filter) TableName() string { return tableFilter }

func (f Filter) BeforeCreate(tx *gorm.DB) error {
	return repositories.AddOnConflictDoNothing(tx, colFilterName, colFilterTypeID)
}

func (f Filter) WhereOnConflict(tx *gorm.DB) *gorm.DB {
	return repositories.WhereOnConflict(tx,
		repositories.ConflictScope{ColumnName: colFilterName, ColumnValue: f.Name.String},
		repositories.ConflictScope{ColumnName: colFilterTypeID, ColumnValue: f.FilterTypeId})
}
