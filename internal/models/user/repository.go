package user

import (
	"fmt"
	"sync"
)

type Repository struct {
	users  map[string]User
	rwlock sync.RWMutex
}

func (r *Repository) Save(user User) {
	r.rwlock.Lock()
	defer r.rwlock.Unlock()

	if r.users == nil {
		r.users = make(map[string]User)
	}
	r.users[user.Name] = user
}

func (r *Repository) FindByUsername(name string) (user User, err error) {
	r.rwlock.RLock()
	defer r.rwlock.RUnlock()

	user, ok := r.users[name]

	if !ok {
		return user, fmt.Errorf("fail to find user")
	}

	return
}
