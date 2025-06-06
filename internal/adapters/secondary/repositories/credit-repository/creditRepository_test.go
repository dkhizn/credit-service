package creditrepository_test

import (
	"testing"

	creditrepository "credit-service/internal/adapters/secondary/repositories/credit-repository"
	credit "credit-service/internal/domain"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreditRepository(t *testing.T) {
	repo := creditrepository.New()
	credit := credit.Credit{
		ObjectCost: decimal.NewFromInt(100000),
	}

	repo.AddCredit(credit)
	cache := repo.GetCache()

	assert.Len(t, cache, 1)
	assert.Equal(t, int64(100000), cache[0].ObjectCost.IntPart())
}
