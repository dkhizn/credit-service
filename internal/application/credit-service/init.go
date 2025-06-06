// Package creditservice реализует бизнес-логику кредитного сервиса.
package creditservice

import (
	creditrepository "credit-service/internal/adapters/secondary/repositories/credit-repository"
	credit "credit-service/internal/domain"

	"github.com/shopspring/decimal"
)

// CreditService определяет интерфейс кредитного сервиса.
type CreditService interface {
	Execute(credit *credit.Credit)
	GetAll() []credit.Credit
}

// CredtiService реализует бизнес-логику кредитного сервиса.
type CredtiService struct {
	rates map[string]decimal.Decimal
	cache *creditrepository.CreditRepository
}

// NewCreditService создаёт новый экземпляр кредитного сервиса.
func NewCreditService() *CredtiService {
	var programRates = map[string]decimal.Decimal{
		"salary":   decimal.NewFromFloat(0.08),
		"military": decimal.NewFromFloat(0.09),
		"base":     decimal.NewFromFloat(0.10),
	}
	return &CredtiService{
		rates: programRates,
		cache: creditrepository.New(),
	}
}
