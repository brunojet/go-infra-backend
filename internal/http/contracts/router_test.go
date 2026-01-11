package contracts

import "testing"

type dummyRouter struct{}

func (dummyRouter) GET(path string, handler HandlerFunc)    {}
func (dummyRouter) POST(path string, handler HandlerFunc)   {}
func (dummyRouter) PATCH(path string, handler HandlerFunc)  {}
func (dummyRouter) DELETE(path string, handler HandlerFunc) {}
func (dummyRouter) Group(path string) Router                { return dummyRouter{} }

func TestRouterInterface(t *testing.T) {
	var _ Router = dummyRouter{}
}
