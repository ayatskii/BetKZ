package main

import (
	"log"
	"net/http"
	"time"

	"BetKZ/internal/api/handlers"
	"BetKZ/internal/repository/memory"
	"BetKZ/internal/service"
)

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

	http.HandleFunc("/users", handlers.CreateUser(userService))
	http.HandleFunc("/events", handlers.ListEvents(eventService))
	http.HandleFunc("/bets", handlers.PlaceBet(betService))

	log.Println("BetKZ server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
