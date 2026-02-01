package domain

type Event struct {
	ID   int64   `json:"id"`
	Name string  `json:"name"`
	Odds float64 `json:"odds"`
	Open bool    `json:"open"`
}
