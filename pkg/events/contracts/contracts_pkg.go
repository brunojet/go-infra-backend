package contracts

import (
	"context"
	"time"
)

// EventSource identifies where the event originated from (filesystem, aws/s3, etc).
type EventSource string

const (
	SourceFilesystem EventSource = "filesystem"
	SourceAWS        EventSource = "aws"
)

// Event represents a generic event payload. Adapters should populate
// Source, Type and Metadata according to the concrete provider (e.g. S3).
type Event struct {
	ID        string         `json:"id,omitempty"`
	Source    EventSource    `json:"source"`
	Type      string         `json:"type"` // e.g. "object.created"
	Timestamp time.Time      `json:"timestamp"`
	Metadata  map[string]any `json:"metadata,omitempty"`
	Payload   any            `json:"payload,omitempty"`
}

// MessageAdapter defines the minimal interface for event adapters. Implementations
// must be able to start/stop and provide a stream of events. Publish is
// optional (not all adapters support publishing).
type MessageAdapter interface {
	// Start begins background processing and returns when ready or on error.
	Start(ctx context.Context) error

	// Stop requests the adapter to stop and waits for cleanup.
	Stop(ctx context.Context) error

	// WaitForMessages performs a blocking long-poll and returns one or more
	// Message envelopes for processing. Manager expects adapters to implement
	// this method for message-style processing. Adapters that cannot support
	// message semantics should return an error.
	WaitForMessages(ctx context.Context) ([]Message, error)
}

// Message represents a received envelope that can be Acked or Nacked.
type Message interface {
	// Event returns the underlying Event payload.
	Event() Event

	// Ack acknowledges successful processing (e.g., DeleteMessage in SQS).
	Ack(ctx context.Context) error

	// Nack signals failed processing. For SQS this could translate to
	// ChangeMessageVisibility (0) or leaving it to expire. Delay hints a
	// requeue delay when supported.
	Nack(ctx context.Context, delay time.Duration) error

	// ExtendVisibility extends the message visibility (heartbeat).
	ExtendVisibility(ctx context.Context, d time.Duration) error
}

// Adapter must provide a blocking long-poll to receive messages from the
// backing system. Manager expects adapters to implement WaitForMessages for
// message-style processing.
// Add the method to Adapter to standardize the interface.
// WaitForMessages should return one or more Message envelopes.
// NOTE: adapters that cannot support message semantics should return an error
// from WaitForMessages.

// WaitForMessages is part of the Adapter contract (see Adapter). Kept here
// as documentation for the required signature.
// func (Adapter) WaitForMessages(ctx context.Context) ([]Message, error)

// Optional convenience type: a function-based handler subscription.
type Handler func(ctx context.Context, evt Event) error
