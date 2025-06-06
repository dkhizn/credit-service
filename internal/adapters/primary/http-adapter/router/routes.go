package router

import (
	"credit-service/internal/adapters/primary/http-adapter/controller"
	"net/http"
)

// RegisterRoutes регистрирует маршруты для контроллера.
func (r *Router) RegisterRoutes(h controller.Handler) {
	routes := []Route{
		{
			Name:    "Execute",
			Method:  http.MethodPost,
			Path:    "/execute",
			Handler: h.Execute,
		},
		{
			Name:    "Cache",
			Method:  http.MethodGet,
			Path:    "/cache",
			Handler: h.Cache,
		},
	}

	for _, rt := range routes {
		r.router.
			Methods(rt.Method).
			Path(rt.Path).
			Name(rt.Name).
			HandlerFunc(rt.Handler)
	}
}
