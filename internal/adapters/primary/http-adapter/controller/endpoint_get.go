package controller

import (
	"encoding/json"
	"log"
	"net/http"
)

// Cache обрабатывает GET-запрос на получение закэшированных данных.
// revive:disable:unused-parameter
// r нужен для сигнатуры http.HandlerFunc.
func (ctr *Controller) Cache(w http.ResponseWriter, r *http.Request) {
	cache := ctr.creditService.GetAll()

	if len(cache) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]string{"error": "empty cache"}); err != nil {
			log.Println("failed to encode empty cache response:", err)
		}

		return
	}

	resp := make([]CacheDTO, 0, 2)
	for i, v := range cache {
		entry := toGetCacheResponseDTO(i, v)
		resp = append(resp, entry)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
