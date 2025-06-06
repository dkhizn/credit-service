// Package router содержит маршруты и роутер для HTTP-сервера.
package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Router конфигурирует маршруты HTTP-сервера.
type Router struct {
	router *mux.Router
}

// NewRouter создаёт новый экземпляр роутера.
func NewRouter() *Router {
	return &Router{
		router: mux.NewRouter(),
	}
}

// Route описывает отдельный маршрут.
type Route struct {
	Name    string
	Method  string
	Path    string
	Handler http.HandlerFunc
}

// Router возвращает настроенный http.Handler с зарегистрированными маршрутами.
func (r *Router) Router() http.Handler {
	return r.router
}
