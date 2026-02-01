package main

import (
	"log"
	"net/http"
	"time"

	"BetKZ/internal/api/handlers"
	"BetKZ/internal/repository/memory"
	"BetKZ/internal/service"
)

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow frontend
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func main() {
	userRepo := memory.NewUserRepository()
	eventRepo := memory.NewEventRepository()
	betRepo := memory.NewBetRepository()

	userService := service.NewUserService(userRepo)
	eventService := service.NewEventService(eventRepo)
	betService := service.NewBetService(betRepo, userRepo, eventRepo)

	// Background worker (goroutine)
	go func() {
		for {
			time.Sleep(10 * time.Second)
			betService.SettleBets()
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/users", handlers.CreateUser(userService))
	mux.HandleFunc("/events", handlers.ListEvents(eventService))
	mux.HandleFunc("/bets", handlers.PlaceBet(betService))

	log.Println("BetKZ server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", withCORS(mux)))
}
