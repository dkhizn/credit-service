package controller

import (
	"bytes"
	credit "credit-service/internal/domain"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

type mockCreditService struct {
	executeFunc func(credit *credit.Credit)
	getAllFunc  func() []credit.Credit
}

func (m *mockCreditService) Execute(credit *credit.Credit) {
	if m.executeFunc != nil {
		m.executeFunc(credit)
	}
}

func (m *mockCreditService) GetAll() []credit.Credit {
	if m.getAllFunc != nil {
		return m.getAllFunc()
	}
	return nil
}

func TestController_Execute_Success(t *testing.T) {
	mockSvc := &mockCreditService{
		executeFunc: func(credit *credit.Credit) {},
	}
	ctrl := New(mockSvc)

	reqBody := ExecuteRequestDTO{
		ObjectCost:     100000,
		InitialPayment: 20000,
		Months:         12,
		Program:        map[string]bool{"base": true},
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/execute", bytes.NewReader(body))
	w := httptest.NewRecorder()

	ctrl.Execute(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestController_Execute_ValidationError(t *testing.T) {
	mockSvc := &mockCreditService{}
	ctrl := New(mockSvc)

	testCases := []struct {
		name     string
		reqBody  ExecuteRequestDTO
		errorMsg string
	}{
		{
			"No program",
			ExecuteRequestDTO{Program: map[string]bool{}},
			"choose program",
		},
		{
			"Multiple programs",
			ExecuteRequestDTO{Program: map[string]bool{"base": true, "salary": true}},
			"choose only 1 program",
		},
		{
			"Low initial payment",
			ExecuteRequestDTO{
				ObjectCost:     100000,
				InitialPayment: 10000,
				Program:        map[string]bool{"base": true},
			},
			"the initial payment should be more",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.reqBody)
			req := httptest.NewRequest("POST", "/execute", bytes.NewReader(body))
			w := httptest.NewRecorder()

			ctrl.Execute(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			var resp map[string]string
			json.Unmarshal(w.Body.Bytes(), &resp)
			assert.Equal(t, tc.errorMsg, resp["error"])
		})
	}
}

func TestToDomain(t *testing.T) {
	dto := ExecuteRequestDTO{
		ObjectCost:     100000,
		InitialPayment: 20000,
		Months:         12,
		Program:        map[string]bool{"salary": true},
	}

	credit := ToDomain(dto)

	assert.Equal(t, decimal.NewFromInt(100000), credit.ObjectCost)
	assert.Equal(t, decimal.NewFromInt(20000), credit.InitialPayment)
	assert.Equal(t, 12, credit.Months)
	assert.Equal(t, "salary", credit.Program)
}

func TestToGetExecuteResponseDTO(t *testing.T) {
	credit := credit.Credit{
		ObjectCost:      decimal.NewFromInt(100000),
		InitialPayment:  decimal.NewFromInt(20000),
		Months:          12,
		Program:         "base",
		RatePercent:     decimal.NewFromFloat(10),
		LoanSum:         decimal.NewFromInt(80000),
		MonthlyPayment:  decimal.NewFromInt(7000),
		Overpayment:     decimal.NewFromInt(4000),
		LastPaymentDate: "2023-12-31",
	}

	resp := toGetExecuteResponseDTO(credit)

	assert.Equal(t, 100000, resp.Result.Params.ObjectCost)
	assert.Equal(t, 80000, resp.Result.Aggregates.LoanSum)
	assert.Equal(t, 10, resp.Result.Aggregates.Rate)
	assert.Equal(t, "2023-12-31", resp.Result.Aggregates.LastPaymentDate)
}

func TestController_Cache_Success(t *testing.T) {
	mockSvc := &mockCreditService{
		getAllFunc: func() []credit.Credit {
			return []credit.Credit{
				{
					ObjectCost:      decimal.NewFromInt(100000),
					InitialPayment:  decimal.NewFromInt(20000),
					Months:          12,
					Program:         "base",
					RatePercent:     decimal.NewFromFloat(10),
					LoanSum:         decimal.NewFromInt(80000),
					MonthlyPayment:  decimal.NewFromInt(7000),
					Overpayment:     decimal.NewFromInt(4000),
					LastPaymentDate: "2023-12-31",
				},
			}
		},
	}
	ctrl := New(mockSvc)

	req := httptest.NewRequest("GET", "/cache", nil)
	w := httptest.NewRecorder()

	ctrl.Cache(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []CacheDTO
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, 100000, resp[0].Params.ObjectCost)
	assert.Equal(t, "2023-12-31", resp[0].Aggregates.LastPaymentDate)
}

func TestController_Cache_Empty(t *testing.T) {
	mockSvc := &mockCreditService{
		getAllFunc: func() []credit.Credit {
			return []credit.Credit{}
		},
	}
	ctrl := New(mockSvc)

	req := httptest.NewRequest("GET", "/cache", nil)
	w := httptest.NewRecorder()

	ctrl.Cache(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "empty cache", resp["error"])
}

func TestController_Execute_InvalidJSON(t *testing.T) {
	mockSvc := &mockCreditService{}
	ctrl := New(mockSvc)

	req := httptest.NewRequest("POST", "/execute", bytes.NewBufferString("{invalid json}"))
	w := httptest.NewRecorder()

	ctrl.Execute(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "invalid request", resp["error"])
}
