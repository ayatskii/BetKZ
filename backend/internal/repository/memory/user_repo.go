package memory

import (
	"sync"

	"BetKZ/internal/domain"
)

type UserRepository struct {
	mu    sync.Mutex
	users map[int64]domain.User
	next  int64
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: make(map[int64]domain.User),
		next:  1,
	}
}

func (r *UserRepository) Create(user domain.User) domain.User {
	r.mu.Lock()
	defer r.mu.Unlock()

	user.ID = r.next
	r.next++
	r.users[user.ID] = user
	return user
}

func (r *UserRepository) GetByID(id int64) (domain.User, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	u, ok := r.users[id]
	return u, ok
}

func (r *UserRepository) Update(user domain.User) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
}
