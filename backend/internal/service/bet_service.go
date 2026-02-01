package service

import "BetKZ/internal/domain"

type BetService struct {
	betRepo interface {
		Create(domain.Bet) domain.Bet
		List() []domain.Bet
		Update(domain.Bet)
	}
	userRepo interface {
		GetByID(int64) (domain.User, bool)
		Update(domain.User)
	}
	eventRepo interface {
		GetByID(int64) (domain.Event, bool)
	}
}

func NewBetService(b, u, e interface{}) *BetService {
	return &BetService{
		betRepo: b.(interface {
			Create(domain.Bet) domain.Bet
			List() []domain.Bet
			Update(domain.Bet)
		}),
		userRepo: u.(interface {
			GetByID(int64) (domain.User, bool)
			Update(domain.User)
		}),
		eventRepo: e.(interface {
			GetByID(int64) (domain.Event, bool)
		}),
	}
}

func (s *BetService) PlaceBet(userID, eventID int64, amount float64) (domain.Bet, bool) {
	user, ok := s.userRepo.GetByID(userID)
	if !ok || user.Balance < amount {
		return domain.Bet{}, false
	}

	event, ok := s.eventRepo.GetByID(eventID)
	if !ok || !event.Open {
		return domain.Bet{}, false
	}

	user.Balance -= amount
	s.userRepo.Update(user)

	return s.betRepo.Create(domain.Bet{
		UserID:  userID,
		EventID: eventID,
		Amount:  amount,
		Status:  "PENDING",
	}), true
}

func (s *BetService) SettleBets() {
	for _, bet := range s.betRepo.List() {
		if bet.Status == "PENDING" {
			bet.Status = "WON"
			s.betRepo.Update(bet)
		}
	}
}
