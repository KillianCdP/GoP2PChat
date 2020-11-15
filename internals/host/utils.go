package host

import (
	"GoP2PChat/internals/user"
	"errors"
)

func getUser(users []user.User, name string) (user.User, error) {
	for _, u := range users {
		if u.Name == name {
			return u, nil
		}
	}

	return user.User{}, errors.New("user not found")
}