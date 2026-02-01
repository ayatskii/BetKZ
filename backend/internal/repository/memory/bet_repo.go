package memory

import (
	"sync"

	"BetKZ/internal/domain"
)

type BetRepository struct {
	mu   sync.Mutex
	bets map[int64]domain.Bet
	next int64
}

func NewBetRepository() *BetRepository {
	return &BetRepository{
		bets: make(map[int64]domain.Bet),
		next: 1,
	}
}

func (r *BetRepository) Create(bet domain.Bet) domain.Bet {
	r.mu.Lock()
	defer r.mu.Unlock()

	bet.ID = r.next
	r.next++
	r.bets[bet.ID] = bet
	return bet
}

func (r *BetRepository) List() []domain.Bet {
	r.mu.Lock()
	defer r.mu.Unlock()

	var result []domain.Bet
	for _, b := range r.bets {
		result = append(result, b)
	}
	return result
}

func (r *BetRepository) Update(bet domain.Bet) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.bets[bet.ID] = bet
}
