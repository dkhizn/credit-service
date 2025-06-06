package creditservice

import (
	credit "credit-service/internal/domain"
)

// Execute выполняет расчёт всех параметров кредита.
func (service *CredtiService) Execute(credit *credit.Credit) {
	credit.Rate = service.rate(credit.Program)
	credit.RatePercent = service.RatePercent(credit.Rate)
	credit.LoanSum = service.loanSum(credit.ObjectCost, credit.InitialPayment)
	credit.MonthlyPayment = service.monthlyPayment(credit.LoanSum, credit.Rate, credit.Months)
	credit.Overpayment = service.overpayment(credit.MonthlyPayment, credit.Months, credit.LoanSum)
	credit.LastPaymentDate = service.LastPaymentDate(credit.Months)

	service.cache.AddCredit(*credit)
}

// GetAll возвращает все доступные кредиты.
func (service *CredtiService) GetAll() []credit.Credit {
	return service.cache.GetCache()
}
