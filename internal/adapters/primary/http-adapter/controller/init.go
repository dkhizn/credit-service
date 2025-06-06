package controller

import (
	creditService "credit-service/internal/application/credit-service"
	"net/http"
)

type Handler interface {
	Execute(w http.ResponseWriter, r *http.Request)
	Cache(w http.ResponseWriter, r *http.Request)
}

// Controller управляет HTTP-запросами, связанными с кредитами.
type Controller struct {
	creditService creditService.CreditService
}

// New создаёт новый контроллер кредитного сервиса.
func New(service creditService.CreditService) *Controller {
	return &Controller{creditService: service}
}
