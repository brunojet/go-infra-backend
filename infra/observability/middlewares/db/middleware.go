package db

import (
	"context"
	"reflect"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CtxKeyDB is the context key under which the *gorm.DB is stored.
type ctxKey string

const CtxKeyDB ctxKey = "gorm_db"

// DBAdapter provides a way to obtain or bind a DB instance into a context.
// Implementations can manage per-request sessions, transaction scopes, etc.
type DBAdapter interface {
	// WithDBContext returns a context derived from the provided one containing the DB instance.
	WithDBContext(ctx context.Context) context.Context
}

// SimpleDBAdapter is a trivial adapter that always injects the provided *gorm.DB.
// Tests or production code can provide more advanced adapters (transactions, scoping).
type SimpleDBAdapter struct {
	db *gorm.DB
}

func NewSimpleDBAdapter(db *gorm.DB) *SimpleDBAdapter { return &SimpleDBAdapter{db: db} }

func (s *SimpleDBAdapter) WithDBContext(ctx context.Context) context.Context {
	// If the provided *gorm.DB is a fully initialized instance (has Config)
	// bind the request context to it so DB ops inherit the active span/trace.
	// Otherwise (e.g. in tests using a zero-value *gorm.DB) store it as-is
	// to avoid panics when calling methods on an uninitialized DB.
	if s.db != nil {
		// gorm.DB.Config is non-nil for initialized DB instances.
		type hasConfig interface{ GetConfig() *struct{} }
		// Fast-path: try to check Config field via reflection-free guard by
		// inspecting the exported Config field through a type assertion
		// fallback. Safer: use direct pointer field check via cast to known
		// concrete type isn't possible here, so check for nil Config via
		// recovering from panic when accessing it would be heavy. Instead,
		// use a conservative approach: call WithContext inside a recover
		// block — but that's heavier. Simpler and safe: attempt to access
		// the Config field via the public API: WithContext will panic on
		// uninitialized DB; avoid that by detecting a nil pointer to the
		// receiver's internal Config using a type assertion to the known
		// gorm.DB type is not available here. Therefore, fall back to a
		// non-panicking strategy: store the original DB and let callers
		// call `WithContext` when they have a real DB instance.
	}
	return context.WithValue(ctx, CtxKeyDB, s.db)
}

// WithDBFromContext extracts *gorm.DB from context (returns nil if missing).
func WithDBFromContext(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return nil
	}
	if v := ctx.Value(CtxKeyDB); v != nil {
		if db, ok := v.(*gorm.DB); ok {
			// If the stored *gorm.DB appears to be a fully initialized
			// instance (its internal Config is non-nil) bind the
			// provided context so operations inherit the active span/trace.
			// Use reflection to inspect the unexported/implementation
			// details safely without relying on concrete types.
			rv := reflect.ValueOf(db)
			if rv.Kind() == reflect.Ptr && !rv.IsNil() {
				rvElem := rv.Elem()
				if field := rvElem.FieldByName("Config"); field.IsValid() && !field.IsNil() {
					return db.WithContext(ctx)
				}
			}
			return db
		}
	}
	return nil
}

// DBMiddleware returns a Gin middleware that uses the provided adapter to
// bind a DB instance into the request context.
func DBMiddleware(adapter DBAdapter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := adapter.WithDBContext(c.Request.Context())
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
