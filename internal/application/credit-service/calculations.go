package creditservice

import (
	"time"

	"github.com/shopspring/decimal"
)

func (s *CredtiService) rate(rate string) decimal.Decimal {
	v, ok := s.rates[rate]
	if !ok {
		return s.rates["base"]
	}

	return v
}

// RatePercent рассчитывает процентную ставку в виде процента.
func (s *CredtiService) RatePercent(rate decimal.Decimal) decimal.Decimal {
	return rate.Mul(decimal.NewFromInt(100))
}

func (s *CredtiService) loanSum(cost, payment decimal.Decimal) decimal.Decimal {
	return cost.Sub(payment)
}

func (s *CredtiService) monthlyRate(rate decimal.Decimal) decimal.Decimal {
	return rate.Div(decimal.NewFromInt(12))
}

func (s *CredtiService) monthlyPayment(loanSum decimal.Decimal, rate decimal.Decimal, months int) decimal.Decimal {
	S := loanSum
	G := s.monthlyRate(rate)
	T := decimal.NewFromInt(int64(months))

	one := decimal.NewFromInt(1)
	powG := one.Add(G).Pow(T)
	numerator := G.Mul(powG)
	numerator = S.Mul(numerator)
	denominator := powG.Sub(one)
	PM := numerator.Div(denominator)
	return PM
}
func (s *CredtiService) overpayment(monthlyPayment decimal.Decimal, months int, loanSum decimal.Decimal) decimal.Decimal {
	monthlyPayment = monthlyPayment.Round(0)
	totalPaid := monthlyPayment.Mul(decimal.NewFromInt(int64(months)))
	return totalPaid.Sub(loanSum)
}

// LastPaymentDate возвращает дату последнего платежа.
func (s *CredtiService) LastPaymentDate(months int) string {
	lastDate := time.Now().AddDate(0, months, 0).Format("2006-01-02")
	return lastDate
}
