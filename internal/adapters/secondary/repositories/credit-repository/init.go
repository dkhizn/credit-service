// Package creditrepository предоставляет доступ к данным о кредитах.
package creditrepository

import (
	credit "credit-service/internal/domain"
	"sync"
)

// CreditRepository реализует хранение и извлечение кредитов из кэша.
type CreditRepository struct {
	data []credit.Credit
	mu   *sync.RWMutex
}

// New создаёт новый репозиторий кредитов.
func New() *CreditRepository {
	return &CreditRepository{
		data: make([]credit.Credit, 0),
		mu:   &sync.RWMutex{},
	}
}
