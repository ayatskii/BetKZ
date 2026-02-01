package domain

type Bet struct {
	ID      int64   `json:"id"`
	UserID  int64   `json:"user_id"`
	EventID int64   `json:"event_id"`
	Amount  float64 `json:"amount"`
	Status  string  `json:"status"`
}
