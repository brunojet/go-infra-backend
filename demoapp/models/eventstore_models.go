package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"reflect"
	"sync"
	"time"

	"gorm.io/gorm"
)

// ESAggregateType mirrors AuditEntityType but kept separate for clarity.
type ESAggregateType int16

const (
	ESAggregateUnknown                    ESAggregateType = 0
	ESAggregateApplication                ESAggregateType = 1
	ESAggregateTerminalModel              ESAggregateType = 2
	ESAggregateTerminalModelConfiguration ESAggregateType = 3
)

// Event is the canonical event-store row. It's optimized for efficient
// per-aggregate reads (ordered by Sequence) and for compact storage of
// event payloads (JSON) and metadata.
type Event struct {
	EventId       int64           `gorm:"primaryKey;autoIncrement"`
	AggregateType ESAggregateType `gorm:"not null;index:idx_es_aggregate_type,priority:1"`

	// AggregateId is stored as binary to allow UUIDs or compact keys.
	AggregateId []byte `gorm:"type:binary(36);not null;uniqueIndex:ux_es_aggregate_seq,priority:1;index:idx_es_aggregate,priority:1"`

	// Sequence is the per-aggregate increasing sequence number.
	Sequence uint64 `gorm:"not null;uniqueIndex:ux_es_aggregate_seq,priority:2;index:idx_es_sequence,priority:2"`

	// EventType is a business-level name for the event (e.g. APP_VERSION_PROMOTED).
	EventType string `gorm:"size:128;not null;index:idx_es_event_type"`

	// Payload contains the typed event payload as JSON.
	Payload json.RawMessage `gorm:"type:json"`

	// Metadata may include correlation id, actor, transport hints, etc.
	Metadata json.RawMessage `gorm:"type:json"`

	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_es_created_at"`
}

func (Event) TableName() string { return "es_event" }

func (Event) WhereOnConflict(tx *gorm.DB) *gorm.DB {
	return tx
}

// Snapshot holds periodic snapshots for fast aggregate reconstruction.
type Snapshot struct {
	SnapshotId    int64           `gorm:"primaryKey;autoIncrement"`
	AggregateType ESAggregateType `gorm:"not null;index:idx_es_snapshot_agg_type,priority:1"`
	AggregateId   []byte          `gorm:"type:binary(36);not null;index:idx_es_snapshot_agg,priority:1"`
	Version       uint64          `gorm:"not null;index:idx_es_snapshot_version,priority:2"`
	Data          json.RawMessage `gorm:"type:json"`
	CreatedAt     time.Time       `gorm:"autoCreateTime"`
}

func (Snapshot) TableName() string { return "es_snapshot" }

func (Snapshot) WhereOnConflict(tx *gorm.DB) *gorm.DB {
	return tx
}

// AppendEvent appends an event for the given aggregate. It computes the next
// sequence number in a transaction to ensure per-aggregate ordering.
func AppendEvent(db *gorm.DB, aggType ESAggregateType, aggId []byte, eventType string, payload interface{}, metadata interface{}) (*Event, error) {
	// marshal payload and metadata
	payloadB, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	metaB, err := json.Marshal(metadata)
	if err != nil {
		// metadata is optional — fall back to null JSON
		metaB = nil
	}

	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	// get current max sequence for this aggregate
	var last sql.NullInt64
	row := tx.Model(&Event{}).Select("COALESCE(MAX(sequence),0)").Where("aggregate_id = ?", aggId).Row()
	if err := row.Scan(&last); err != nil {
		tx.Rollback()
		return nil, err
	}
	seq := uint64(last.Int64) + 1

	ev := &Event{
		AggregateType: aggType,
		AggregateId:   aggId,
		Sequence:      seq,
		EventType:     eventType,
		Payload:       json.RawMessage(payloadB),
		Metadata:      json.RawMessage(metaB),
	}

	if err := tx.Create(ev).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return ev, nil
}

// GetEvents returns events for an aggregate ordered by Sequence ascending.
func GetEvents(db *gorm.DB, aggId []byte, fromSequence uint64, limit int) ([]Event, error) {
	var evs []Event
	q := db.Model(&Event{}).Where("aggregate_id = ? AND sequence >= ?", aggId, fromSequence).Order("sequence ASC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if err := q.Find(&evs).Error; err != nil {
		return nil, err
	}
	return evs, nil
}

// EncodeEntityPK canonicalizes an aggregate PK (struct, map, string, number) into
// a fixed-size byte array suitable for indexing. Strategy: JSON marshal (caller
// may provide a canonical map/struct) then sha256 truncated to 18 bytes.
func EncodeEntityPK(pk interface{}) ([]byte, error) {
	if pk == nil {
		return nil, nil
	}
	// Fast-path: if already []byte of reasonable length, return padded/truncated
	if b, ok := pk.([]byte); ok {
		if len(b) <= 36 {
			out := make([]byte, 36)
			copy(out, b)
			return out, nil
		}
		// truncate
		out := make([]byte, 36)
		copy(out, b[:36])
		return out, nil
	}

	// Marshal to JSON then hash/truncate
	jb, err := json.Marshal(pk)
	if err != nil {
		return nil, err
	}
	h := sha256.Sum256(jb)
	out := make([]byte, 18)
	copy(out, h[:18])
	return out, nil
}

// DTO mapper registry: allows the application to register per-DTO mappers that
// produce typed event type, payload and metadata. If no mapper exists, the DTO
// will be marshaled as JSON and used as payload with an inferred event type.
var (
	dtoMappersMu sync.RWMutex
	dtoMappers   = map[string]func(interface{}) (eventType string, payload interface{}, metadata interface{}, err error){}
)

// RegisterDTOMapper registers a mapper function for a DTO type. Use
// reflect.TypeOf(exampleDTO).String() as key when registering.
func RegisterDTOMapper(exampleDTO interface{}, mapper func(interface{}) (string, interface{}, interface{}, error)) {
	if exampleDTO == nil || mapper == nil {
		return
	}
	key := reflect.TypeOf(exampleDTO).String()
	dtoMappersMu.Lock()
	dtoMappers[key] = mapper
	dtoMappersMu.Unlock()
}

// AppendTypedEvent is a convenience wrapper that encodes aggregate PK,
// maps DTO->event (using registered mapper) or falls back to JSON payload,
// and appends the event.
func AppendTypedEvent(db *gorm.DB, aggType ESAggregateType, aggPk interface{}, dto interface{}, actor string) (*Event, error) {
	aggId, err := EncodeEntityPK(aggPk)
	if err != nil {
		return nil, err
	}

	// Try mapper
	var eventType string
	var payload interface{}
	var metadata interface{}

	if dto != nil {
		key := reflect.TypeOf(dto).String()
		dtoMappersMu.RLock()
		mapper, ok := dtoMappers[key]
		dtoMappersMu.RUnlock()
		if ok {
			et, p, m, err := mapper(dto)
			if err != nil {
				return nil, err
			}
			eventType = et
			payload = p
			metadata = m
		} else {
			// fallback: use DTO as payload directly (will be marshaled to JSON by AppendEvent)
			eventType = reflect.TypeOf(dto).Name()
			payload = dto
			metadata = map[string]string{"actor": actor}
		}
	} else {
		eventType = "UNKNOWN"
		payload = nil
		metadata = map[string]string{"actor": actor}
	}

	return AppendEvent(db, aggType, aggId, eventType, payload, metadata)
}
