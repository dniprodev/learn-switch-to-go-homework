package handlers

import (
	"github.com/dniprodev/learn-switch-to-go-homework/internal/models/user"
)

type UserRepositoryInterface interface {
	FindByUsername(username string) (user.User, bool)
	Save(user user.User) user.User
}
