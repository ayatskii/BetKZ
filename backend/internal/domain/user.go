package domain

type User struct {
	ID      int64   `json:"id"`
	Email   string  `json:"email"`
	Balance float64 `json:"balance"`
}
