package creditrepository

import (
	credit "credit-service/internal/domain"
)

// AddCredit добавляет кредит в кэш.
func (rep *CreditRepository) AddCredit(credit credit.Credit) {
	rep.mu.Lock()
	rep.data = append(rep.data, credit)
	rep.mu.Unlock()
}

// GetCache возвращает закэшированные кредиты.
func (rep *CreditRepository) GetCache() []credit.Credit {
	copySlice := make([]credit.Credit, len(rep.data))
	rep.mu.Lock()
	copy(copySlice, rep.data)
	rep.mu.Unlock()

	return copySlice
}
