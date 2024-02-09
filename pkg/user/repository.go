package user

type Repository struct {
	users map[string]User
}

func (r *Repository) save(user User) {
	if r.users == nil {
		r.users = make(map[string]User)
	}
	r.users[user.name] = user
}

func (r *Repository) findByUsername(name string) (user User, ok bool) {
	user, ok = r.users[name]
	return
}
