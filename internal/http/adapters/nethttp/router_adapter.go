package nethttp

import (
	"net/http"

	"github.com/brunojet/go-infra-backend/internal/http/contracts"
	"github.com/go-chi/chi/v5"
)

type ChiRouter struct {
	r chi.Router
}

func NewChiRouter(r chi.Router) contracts.Router {
	return &ChiRouter{r: r}
}

func (cr *ChiRouter) GET(path string, handler contracts.HandlerFunc) {
	cr.r.Get(translatePath(path), func(w http.ResponseWriter, r *http.Request) {
		handler(NewContext(w, r))
	})
}

func (cr *ChiRouter) POST(path string, handler contracts.HandlerFunc) {
	cr.r.Post(translatePath(path), func(w http.ResponseWriter, r *http.Request) {
		handler(NewContext(w, r))
	})
}

func (cr *ChiRouter) PATCH(path string, handler contracts.HandlerFunc) {
	cr.r.Method(http.MethodPatch, translatePath(path), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(NewContext(w, r))
	}))
}

func (cr *ChiRouter) DELETE(path string, handler contracts.HandlerFunc) {
	cr.r.Delete(translatePath(path), func(w http.ResponseWriter, r *http.Request) {
		handler(NewContext(w, r))
	})
}

func (cr *ChiRouter) Group(relativePath string) contracts.Router {
	sub := chi.NewRouter()
	cr.r.Mount(translatePath(relativePath), sub)
	return &ChiRouter{r: sub}
}
