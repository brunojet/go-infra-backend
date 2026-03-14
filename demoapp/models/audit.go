package models

import (
	"database/sql"
	"time"
)

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
	AuditEntityUnknown AuditEntityType = iota
	AuditEntityApplication
	AuditEntityTerminalModel
	AuditEntityTerminalModelConfiguration
	AuditEntityApplicationConfiguration
	AuditEntityFilterType
	AuditEntityFilter
	AuditEntityApplicationProfile
	AuditEntityApplicationVersion
	AuditEntityApplicationCatalog
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
