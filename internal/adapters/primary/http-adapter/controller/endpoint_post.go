package controller

import (
	"encoding/json"
	"log"
	"net/http"
)

// Execute обрабатывает POST-запрос на расчёт кредита.
// revive:disable:unused-parameter
// r нужен для сигнатуры http.HandlerFunc.
func (ctr *Controller) Execute(w http.ResponseWriter, r *http.Request) {
	var executeRequestDTO ExecuteRequestDTO
	err := json.NewDecoder(r.Body).Decode(&executeRequestDTO)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"}); err != nil {
			log.Println("failed to encode invalid request response:", err)
		}
		return
	}

	err = validate(executeRequestDTO)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}); err != nil {
			log.Println("failed to encode invalid request response:", err)
		}

		return
	}

	creditDomain := ToDomain(executeRequestDTO)
	ctr.creditService.Execute(&creditDomain)
	response := toGetExecuteResponseDTO(creditDomain)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
