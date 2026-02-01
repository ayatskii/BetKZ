package handlers

import (
	"encoding/json"
	"net/http"

	"BetKZ/internal/service"
)

func CreateUser(s *service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Email   string  `json:"email"`
			Balance float64 `json:"balance"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		user := s.CreateUser(req.Email, req.Balance)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}
