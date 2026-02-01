package handlers

import (
	"encoding/json"
	"net/http"

	"BetKZ/internal/service"
)

func ListEvents(s *service.EventService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(s.ListEvents())
	}
}
