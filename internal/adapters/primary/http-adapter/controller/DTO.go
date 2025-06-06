// Package controller реализует обработчики HTTP-ручек и содержит DTO-структуры для обмена данными между слоями.
package controller

import (
	credit "credit-service/internal/domain"

	"github.com/shopspring/decimal"
)

// ParamsDTO описывает изначальные входные параметры долга.
type ParamsDTO struct {
	ObjectCost     int `json:"object_cost"`
	InitialPayment int `json:"initial_payment"`
	Months         int `json:"months"`
}

// AggregatesDTO описывает параметры, рассчитываемые ипотечным калькулятором.
type AggregatesDTO struct {
	Rate            int    `json:"rate"`
	LoanSum         int    `json:"loan_sum"`
	MonthlyPayment  int    `json:"monthly_payment"`
	Overpayment     int    `json:"overpayment"`
	LastPaymentDate string `json:"last_payment_date"`
}

// ExecuteResponseDTO описывает параметры, которые отдаем в качестве ответа в /execute.
type ExecuteResponseDTO struct {
	Result struct {
		Params     ParamsDTO       `json:"params"`
		Program    map[string]bool `json:"program"`
		Aggregates AggregatesDTO   `json:"aggregates"`
	} `json:"result"`
}

// ExecuteRequestDTO описывает параметры, которые получаем при /execute.
type ExecuteRequestDTO struct {
	ObjectCost     int             `json:"object_cost"`
	InitialPayment int             `json:"initial_payment"`
	Months         int             `json:"months"`
	Program        map[string]bool `json:"program"`
}

// CacheDTO описывает параметры, которые отдаем в качестве ответа в /cache.
type CacheDTO struct {
	ID     int `json:"id"`
	Params struct {
		ObjectCost     int `json:"object_cost"`
		InitialPayment int `json:"initial_payment"`
		Months         int `json:"months"`
	} `json:"params"`
	Program    map[string]bool `json:"program"`
	Aggregates struct {
		Rate            int    `json:"rate"`
		LoanSum         int    `json:"loan_sum"`
		MonthlyPayment  int    `json:"monthly_payment"`
		Overpayment     int    `json:"overpayment"`
		LastPaymentDate string `json:"last_payment_date"`
	} `json:"aggregates"`
}

// ToDomain возвращает структуру домена, построенную из DTO.
func ToDomain(dto ExecuteRequestDTO) credit.Credit {
	var program string

	switch {
	case dto.Program["salary"]:
		program = "salary"
	case dto.Program["military"]:
		program = "military"
	default:
		program = "base"
	}

	return credit.Credit{
		ObjectCost:     decimal.NewFromInt(int64(dto.ObjectCost)),
		InitialPayment: decimal.NewFromInt(int64(dto.InitialPayment)),
		Months:         dto.Months,
		Program:        program,
	}
}

func toGetExecuteResponseDTO(credit credit.Credit) ExecuteResponseDTO {
	programMap := make(map[string]bool)
	programMap[credit.Program] = true

	return ExecuteResponseDTO{
		Result: struct {
			Params     ParamsDTO       `json:"params"`
			Program    map[string]bool `json:"program"`
			Aggregates AggregatesDTO   `json:"aggregates"`
		}{
			Params: ParamsDTO{
				ObjectCost:     roundToInt(credit.ObjectCost),
				InitialPayment: roundToInt(credit.InitialPayment),
				Months:         credit.Months,
			},
			Program: programMap,
			Aggregates: AggregatesDTO{
				Rate:            roundToInt(credit.RatePercent),
				LoanSum:         roundToInt(credit.LoanSum),
				MonthlyPayment:  roundToInt(credit.MonthlyPayment),
				Overpayment:     roundToInt(credit.Overpayment),
				LastPaymentDate: credit.LastPaymentDate,
			},
		},
	}
}

func toGetCacheResponseDTO(id int, credit credit.Credit) CacheDTO {
	var dto CacheDTO
	program := make(map[string]bool)
	program[credit.Program] = true

	dto.ID = id
	dto.Params.ObjectCost = roundToInt(credit.ObjectCost)
	dto.Params.InitialPayment = roundToInt(credit.InitialPayment)
	dto.Params.Months = credit.Months
	dto.Program = program
	dto.Aggregates.Rate = roundToInt(credit.RatePercent)
	dto.Aggregates.LoanSum = roundToInt(credit.LoanSum)
	dto.Aggregates.MonthlyPayment = roundToInt(credit.MonthlyPayment)
	dto.Aggregates.Overpayment = roundToInt(credit.Overpayment)
	dto.Aggregates.LastPaymentDate = credit.LastPaymentDate

	return dto
}
