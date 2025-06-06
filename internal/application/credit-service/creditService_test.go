package creditservice_test

import (
	"testing"
	"time"

	creditService "credit-service/internal/application/credit-service"
	credit "credit-service/internal/domain"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreditService_Execute(t *testing.T) {
	svc := creditService.NewCreditService()
	credit := &credit.Credit{
		ObjectCost:     decimal.NewFromInt(100000),
		InitialPayment: decimal.NewFromInt(20000),
		Months:         12,
		Program:        "base",
	}

	svc.Execute(credit)

	assert.True(t, credit.LoanSum.Equal(decimal.NewFromInt(80000)))
	assert.True(t, credit.Rate.Equal(decimal.NewFromFloat(0.1)))
	assert.True(t, credit.RatePercent.Equal(decimal.NewFromInt(10)))
	assert.True(t, credit.MonthlyPayment.GreaterThan(decimal.Zero))
	assert.True(t, credit.Overpayment.GreaterThan(decimal.Zero))
}

func TestLastPaymentDate(t *testing.T) {
	svc := creditService.NewCreditService()
	date := svc.LastPaymentDate(12)

	expectedDate := time.Now().AddDate(0, 12, 0).Format("2006-01-02")
	assert.Equal(t, expectedDate, date)
}

func TestRatePercent(t *testing.T) {
	svc := creditService.NewCreditService()
	percent := svc.RatePercent(decimal.NewFromFloat(0.15))
	assert.True(t, percent.Equal(decimal.NewFromInt(15)))
}
