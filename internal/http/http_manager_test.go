package http

import (
	"errors"
	"testing"

	"github.com/brunojet/go-infra-backend/internal/http/contracts"
	"github.com/brunojet/go-infra-backend/internal/http/types"
)

func TestHTTPManager_Close(t *testing.T) {
	mgr := &HTTPManager{}
	mgr.runtime = &HttpRuntime{}
	err := mgr.Close()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if mgr.runtime != nil {
		t.Errorf("expected runtime to be nil after Close, got %v", mgr.runtime)
	}
}

func TestNewHTTPManager_UnsupportedDriver(t *testing.T) {
	params := HTTPParams{Driver: types.HTTPDriver("unknown")}
	_, err := NewHTTPManager(params)
	if err == nil {
		t.Error("expected error for unsupported driver")
	}
}

func TestHTTPManager_Open_RuntimeReuse(t *testing.T) {
	mgr := &HTTPManager{}
	rt := &HttpRuntime{}
	mgr.runtime = rt
	got, err := mgr.Open()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != rt {
		t.Errorf("expected same runtime instance")
	}
}

func TestHTTPManager_buildRuntime_Unsupported(t *testing.T) {
	mgr := &HTTPManager{Params: HTTPParams{Driver: types.HTTPDriver("unknown")}}
	_, err := mgr.buildRuntime()
	if err == nil {
		t.Error("expected error for unsupported driver")
	}
}

func TestHTTPManager_OpenAndRegister(t *testing.T) {
	// Register=false, should not call registrar
	mgr := &HTTPManager{Params: HTTPParams{Register: false}}
	mgr.runtime = &HttpRuntime{}
	called := false
	_, err := mgr.OpenAndRegister(func(r contracts.Router) error { called = true; return nil })
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if called {
		t.Error("registrar should not be called when Register is false")
	}

	// Register=true, registrar is nil
	mgr = &HTTPManager{Params: HTTPParams{Register: true}}
	mgr.runtime = &HttpRuntime{}
	_, err = mgr.OpenAndRegister(nil)
	if err != ErrMissingRegistrar {
		t.Errorf("expected ErrMissingRegistrar, got %v", err)
	}

	// Register=true, registrar returns error
	mgr = &HTTPManager{Params: HTTPParams{Register: true}}
	mgr.runtime = &HttpRuntime{Router: nil}
	regErr := errors.New("fail")
	_, err = mgr.OpenAndRegister(func(r contracts.Router) error { return regErr })
	if err != regErr {
		t.Errorf("expected registrar error, got %v", err)
	}
}

func TestHTTPManager_OpenAndRegister_OpenError(t *testing.T) {
	mgr := &HTTPManager{Params: HTTPParams{Driver: types.HTTPDriver("unknown")}}
	_, err := mgr.OpenAndRegister(func(r contracts.Router) error { return nil })
	expected := "unsupported HTTP_DRIVER \"unknown\""
	if err == nil || err.Error() != expected {
		t.Errorf("expected %q, got %v", expected, err)
	}
}

type DummyRouter struct{}

func (d DummyRouter) DELETE(string, contracts.HandlerFunc) {}
func (d DummyRouter) Group(string) contracts.Router        { return d }
func (d DummyRouter) GET(string, contracts.HandlerFunc)    {}
func (d DummyRouter) POST(string, contracts.HandlerFunc)   {}
func (d DummyRouter) PATCH(string, contracts.HandlerFunc)  {}

func TestHTTPManager_OpenAndRegister_registrarPanic(t *testing.T) {
	mgr := &HTTPManager{Params: HTTPParams{Register: true}}
	mgr.runtime = &HttpRuntime{Router: DummyRouter{}}
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic from registrar, got none")
		}
	}()
	_, _ = mgr.OpenAndRegister(func(r contracts.Router) error {
		panic("fail hard")
	})
}
