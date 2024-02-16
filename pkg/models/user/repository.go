package user

import "github.com/google/uuid"

type Repository struct {
	users map[string]User
}

// TODO: Q - Is it ok to return enriched user?
func (r *Repository) Save(user User) User {
	if r.users == nil {
		r.users = make(map[string]User)
	}
	user.ID = uuid.New().String()
	r.users[user.Name] = user

	return user
}

func (r *Repository) FindByUsername(name string) (user User, ok bool) {
	user, ok = r.users[name]
	return
}
