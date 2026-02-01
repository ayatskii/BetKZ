package handlers

import (
	"encoding/json"
	"net/http"

	"BetKZ/internal/service"
)

func PlaceBet(s *service.BetService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			UserID  int64   `json:"user_id"`
			EventID int64   `json:"event_id"`
			Amount  float64 `json:"amount"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		bet, ok := s.PlaceBet(req.UserID, req.EventID, req.Amount)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(bet)
	}
}
