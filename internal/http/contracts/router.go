package contracts

type Router interface {
	GET(path string, handler HandlerFunc)
	POST(path string, handler HandlerFunc)
	PATCH(path string, handler HandlerFunc)
	DELETE(path string, handler HandlerFunc)

	// Group creates a sub-router under the given relative path.
	// Implementations may panic if grouping is not supported.
	Group(relativePath string) Router
}
