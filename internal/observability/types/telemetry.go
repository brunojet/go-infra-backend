package types

// Field is a simple key/value pair used to enrich logs, spans and metrics.
//
// Values should be low-cardinality whenever possible.
type Field struct {
	Key   string
	Value any
}
