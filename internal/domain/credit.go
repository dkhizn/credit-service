// Package credit содержит доменные сущности, связанные с кредитами.
package credit

import "github.com/shopspring/decimal"

// Credit представляет данные о кредите.
type Credit struct {
	ObjectCost      decimal.Decimal
	InitialPayment  decimal.Decimal
	Months          int
	Program         string
	Rate            decimal.Decimal
	RatePercent     decimal.Decimal
	LoanSum         decimal.Decimal
	MonthlyPayment  decimal.Decimal
	Overpayment     decimal.Decimal
	LastPaymentDate string
}
