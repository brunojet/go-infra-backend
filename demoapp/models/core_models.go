package models

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type TerminalModel struct {
	TerminalModelId             int64                        `gorm:"primaryKey;autoIncrement"`
	Name                        sql.NullString               `gorm:"not null;size:32;uniqueIndex:ux_terminal_model_name"`
	Description                 sql.NullString               `gorm:"size:500"`
	CreatedAt                   sql.NullTime                 `gorm:"autoCreateTime;index:idx_terminal_model_del_created,priority:2"`
	UpdatedAt                   sql.NullTime                 `gorm:"autoUpdateTime;index:idx_terminal_model_del_updated,priority:2"`
	DeletedAt                   gorm.DeletedAt               `gorm:"index:idx_terminal_model_del_created,priority:1;index:idx_terminal_model_del_updated,priority:1"`
	TerminalModelConfigurations []TerminalModelConfiguration `gorm:"foreignKey:TerminalModelId;references:TerminalModelId"`
}

func (TerminalModel) TableName() string { return "terminal_model" }

type TerminalModelConfiguration struct {
	TerminalModelConfigurationId int64                      `gorm:"primaryKey;autoIncrement"`
	TerminalModelId              int64                      `gorm:"not null;uniqueIndex:uk_terminal_model_configuration_terminal_integration,priority:1;index:idx_terminal_model_configuration_terminal,priority:1"`
	IntegrationType              sql.NullInt16              `gorm:"not null;uniqueIndex:uk_terminal_model_configuration_terminal_integration,priority:2"`
	CreatedAt                    sql.NullTime               `gorm:"autoCreateTime;index:idx_terminal_model_configuration_del_created,priority:2"`
	UpdatedAt                    sql.NullTime               `gorm:"autoUpdateTime;index:idx_terminal_model_configuration_del_updated,priority:2"`
	DeletedAt                    gorm.DeletedAt             `gorm:"index:idx_terminal_model_configuration_del_created,priority:1;index:idx_terminal_model_configuration_del_updated,priority:1"`
	ApplicationConfigurations    []ApplicationConfiguration `gorm:"foreignKey:TerminalModelConfigurationId;references:TerminalModelConfigurationId"`
}

func (TerminalModelConfiguration) TableName() string { return "terminal_model_configuration" }

type FilterType struct {
	FilterTypeId int64          `gorm:"primaryKey;autoIncrement"`
	Name         sql.NullString `gorm:"size:128;index:ux_filter_type_name"`
	Description  sql.NullString `gorm:"size:500"`
	CreatedAt    sql.NullTime   `gorm:"autoCreateTime;index:idx_filter_type_del_created,priority:2"`
	UpdatedAt    sql.NullTime   `gorm:"autoUpdateTime;index:idx_filter_type_del_updated,priority:2"`
	DeletedAt    gorm.DeletedAt `gorm:"index:idx_filter_type_del_created,priority:1;index:idx_filter_type_del_updated,priority:1"`
	Filters      []Filter       `gorm:"foreignKey:FilterTypeId;references:FilterTypeId"`
}

func (FilterType) TableName() string { return "filter_type" }

type Filter struct {
	FilterId     int64          `gorm:"primaryKey;autoIncrement"`
	FilterTypeId int64          `gorm:"not null;index:idx_filter_type"`
	Name         sql.NullString `gorm:"size:128;index:ux_filter_name_type,priority:1"`
	Description  sql.NullString `gorm:"size:500"`
	CreatedAt    sql.NullTime   `gorm:"autoCreateTime;index:idx_filter_del_created,priority:2"`
	UpdatedAt    sql.NullTime   `gorm:"autoUpdateTime;index:idx_filter_del_updated,priority:2"`
	DeletedAt    gorm.DeletedAt `gorm:"index:idx_filter_del_created,priority:1;index:idx_filter_del_updated,priority:1"`
}

type Application struct {
	ApplicationId int64          `gorm:"primaryKey;autoIncrement"`
	Name          sql.NullString `gorm:"not null;size:32;uniqueIndex:ux_application_name"`
	Description   sql.NullString `gorm:"size:500"`
	CreatedAt     sql.NullTime   `gorm:"autoCreateTime;index:idx_application_del_created,priority:2"`
	UpdatedAt     sql.NullTime   `gorm:"autoUpdateTime;index:idx_application_del_updated,priority:2"`
	DeletedAt     gorm.DeletedAt `gorm:"index:idx_application_del_created,priority:1;index:idx_application_del_updated,priority:1"`

	ApplicationConfigurations []ApplicationConfiguration `gorm:"foreignKey:ApplicationId;references:ApplicationId"`
	ApplicationProfiles       []ApplicationProfile       `gorm:"foreignKey:ApplicationId;references:ApplicationId"`
	ApplicationCatalogs       []ApplicationCatalog       `gorm:"foreignKey:ApplicationId;references:ApplicationId"`
}

func (Application) TableName() string { return "application" }

type ApplicationConfigurationPk struct {
	ApplicationId                int64 `gorm:"column:application_id;primaryKey;priority:1;index:idx_app_cfg_terminal_app,priority:2"`
	TerminalModelConfigurationId int64 `gorm:"column:terminal_model_configuration_id;primaryKey;priority:2;index:idx_app_cfg_terminal_app,priority:1"`
}

type ApplicationConfiguration struct {
	ApplicationConfigurationPk `gorm:"embedded"`
	PackageName                sql.NullString       `gorm:"column:package_name;not null;size:255"`
	CreatedAt                  sql.NullTime         `gorm:"autoCreateTime;index:idx_application_configuration_del_created,priority:2"`
	UpdatedAt                  sql.NullTime         `gorm:"autoUpdateTime;index:idx_application_configuration_del_updated,priority:2"`
	DeletedAt                  gorm.DeletedAt       `gorm:"index:idx_application_configuration_del_created,priority:1;index:idx_application_configuration_del_updated,priority:1"`
	ApplicationVersions        []ApplicationVersion `gorm:"foreignKey:ApplicationId,TerminalModelConfigurationId;references:ApplicationId,TerminalModelConfigurationId"`
	ApplicationCatalogs        []ApplicationCatalog `gorm:"foreignKey:ApplicationId,TerminalModelConfigurationId;references:ApplicationId,TerminalModelConfigurationId"`
}

func (ApplicationConfiguration) TableName() string { return "application_configuration" }

type ApplicationProfile struct {
	ApplicationProfileId int64          `gorm:"primaryKey;autoIncrement"`
	ApplicationId        int64          `gorm:"index"`
	Name                 sql.NullString `gorm:"size:128;index:ux_application_profile_name_app,priority:1"`
	Description          sql.NullString `gorm:"size:255"`
	ReviewAt             sql.NullTime
	ProductionAt         sql.NullTime
	CreatedAt            sql.NullTime         `gorm:"autoCreateTime;index:idx_application_profile_history_del_created,priority:2"`
	UpdatedAt            sql.NullTime         `gorm:"autoUpdateTime;index:idx_application_profile_history_del_updated,priority:2"`
	DeletedAt            gorm.DeletedAt       `gorm:"index:idx_application_profile_history_del_created,priority:1;index:idx_application_profile_history_del_updated,priority:1"`
	Filters              []Filter             `gorm:"many2many:application_profile_filters;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	ApplicationCatalogs  []ApplicationCatalog `gorm:"foreignKey:ApplicationProfileId;references:ApplicationProfileId"`
}

func (ApplicationProfile) TableName() string {
	return "application_profile_history"
}

type ApplicationVersion struct {
	ApplicationVersionId         int64 `gorm:"primaryKey;autoIncrement"`
	ApplicationId                int64 `gorm:"column:application_id;index:idx_appver_app_cfg,priority:1;index:idx_appver_terminal_app,priority:2"`
	TerminalModelConfigurationId int64 `gorm:"column:terminal_model_configuration_id;index:idx_appver_app_cfg,priority:2;index:idx_appver_terminal_app,priority:1"`
	ReviewAt                     sql.NullTime
	ProductionAt                 sql.NullTime
	CreatedAt                    sql.NullTime   `gorm:"autoCreateTime;index:idx_application_version_history_del_created,priority:2"`
	UpdatedAt                    sql.NullTime   `gorm:"autoUpdateTime;index:idx_application_version_history_del_updated,priority:2"`
	DeletedAt                    gorm.DeletedAt `gorm:"index:idx_application_version_history_del_created,priority:1;index:idx_application_version_history_del_updated,priority:1"`
}

func (ApplicationVersion) TableName() string {
	return "application_version_history"
}

type ApplicationCatalogPk struct {
	ApplicationId                int64 `gorm:"column:application_id;primaryKey;priority:1;index:idx_ctlg_terminal_stage_app,priority:3"`
	TerminalModelConfigurationId int64 `gorm:"column:terminal_model_configuration_id;primaryKey;priority:2;index:idx_ctlg_terminal_stage_app,priority:1"`
	Stage                        int16 `gorm:"column:stage;primaryKey;priority:3;index:idx_ctlg_terminal_stage_app,priority:2"`
}

type ApplicationCatalog struct {
	ApplicationCatalogPk `gorm:"embedded"`
	ApplicationVersionId int64          `gorm:"column:application_version_id"`
	ApplicationProfileId int64          `gorm:"column:application_profile_id"`
	CreatedAt            sql.NullTime   `gorm:"autoCreateTime;index:idx_application_catalog_del_created,priority:2"`
	UpdatedAt            sql.NullTime   `gorm:"autoUpdateTime;index:idx_application_catalog_del_updated,priority:2"`
	DeletedAt            gorm.DeletedAt `gorm:"index:idx_application_catalog_del_created,priority:1;index:idx_application_catalog_del_updated,priority:1"`
}

func (ApplicationCatalog) TableName() string { return "application_catalog" }

// AuditOperation represents a generic CRUD operation recorded in the audit log.
// Use small integers to keep storage and indexes compact across databases.
type AuditOperation int16

const (
	AuditOpUnknown AuditOperation = 0
	AuditOpCreate  AuditOperation = 1
	AuditOpUpdate  AuditOperation = 2
	AuditOpDelete  AuditOperation = 3
)

// AuditEntityType is an application-level "enum" stored as small integer.
// It avoids typos and keeps indexes compact, while allowing the audit log to remain generic.
type AuditEntityType int16

const (
	AuditEntityUnknown                    AuditEntityType = 0
	AuditEntityApplication                AuditEntityType = 1
	AuditEntityTerminalModel              AuditEntityType = 2
	AuditEntityTerminalModelConfiguration AuditEntityType = 3
	AuditEntityApplicationConfiguration   AuditEntityType = 4
	AuditEntityApplicationProfile         AuditEntityType = 5
	AuditEntityApplicationVersion         AuditEntityType = 6
	AuditEntityApplicationCatalog         AuditEntityType = 7
)

// FieldValueType describes the stored field value format in AuditFieldChange.
type FieldValueType int16

const (
	FieldValueTypeUnknown FieldValueType = 0
	FieldValueTypeString  FieldValueType = 1
	FieldValueTypeNumber  FieldValueType = 2
	FieldValueTypeBool    FieldValueType = 3
	FieldValueTypeJSON    FieldValueType = 4
)

// AuditEvent is a generic, append-only audit log entry.
//
// It is intentionally decoupled from domain tables to support any entity,
// including composite primary keys. Use EntityPk as a canonical representation
// (e.g. JSON: {"application_id":1,"terminal_model_configuration_id":2}).
type AuditEvent struct {
	AuditEventId int64 `gorm:"primaryKey;autoIncrement"`

	OccurredAt time.Time `gorm:"not null;autoCreateTime;index:idx_audit_event_occurred_at;index:idx_audit_event_actor_occurred,priority:2;index:idx_audit_entity_lookup,priority:3"`

	// ActorUserId should be an external identifier (e.g. UUID from JWT claims).
	ActorUserId sql.NullString `gorm:"size:64;index:idx_audit_event_actor_occurred,priority:1"`

	// EventName is a business-level name (e.g. APP_VERSION_PROMOTED).
	EventName sql.NullString `gorm:"size:128;index:idx_audit_event_name"`

	// EntityType is a compact identifier for the audited entity.
	EntityType AuditEntityType `gorm:"not null;index:idx_audit_event_entity_type,priority:1;index:idx_audit_entity_lookup,priority:1"`

	Operation AuditOperation `gorm:"not null;index:idx_audit_event_entity,priority:3"`

	// EntityPkBytes is a compact, fixed-size binary representation of the entity PK
	// (canonicalized + encoded/truncated to 18 bytes). Use this for indexed lookups.
	EntityPkBytes []byte `gorm:"type:binary(18);index:idx_audit_entity_lookup,priority:2"`

	// CorrelationId can be used to group multiple audit events in a single request/transaction.
	CorrelationId sql.NullString `gorm:"size:64;index:idx_audit_event_correlation"`

	FieldChanges []AuditFieldChange `gorm:"foreignKey:AuditEventId;references:AuditEventId"`
}

func (AuditEvent) TableName() string { return "audit_event" }

// AuditFieldChange stores per-field changes associated with an AuditEvent.
//
// For INSERT, OldValue is NULL; for DELETE, NewValue is NULL.
// Values are stored as strings to keep this generic; ValueType can help parsing.
type AuditFieldChange struct {
	AuditFieldChangeId int64 `gorm:"primaryKey;autoIncrement"`
	AuditEventId       int64 `gorm:"not null;index:idx_audit_field_event,priority:1"`

	FieldName string `gorm:"not null;size:128;index:idx_audit_field_field,priority:1"`
	// FieldValueType is a compact enum describing the stored value format.
	FieldValueType FieldValueType `gorm:"not null;index:idx_audit_field_value_type,priority:2"`

	// FieldValue stores the value that was written (for inserts/updates).
	// For deletes, this can contain the last known state if desired.
	FieldValue sql.NullString `gorm:"type:text"`

	AuditEvent *AuditEvent `gorm:"foreignKey:AuditEventId;references:AuditEventId"`
}

func (AuditFieldChange) TableName() string { return "audit_field_change" }
