package models

import "github.com/vishenosik/sso/internal/entities"

type User struct {
	Nickname string `db:"nickname"`
	Email    string `db:"email"`
	ID       string `db:"id"`
	Password []byte `db:"password"`
}

func UserToEntities(user *User) *entities.UserCreds {
	if user == nil {
		return new(entities.UserCreds)
	}
	return &entities.UserCreds{
		User: entities.User{
			ID:    user.ID,
			Email: user.Email,
		},
		Password: string(user.Password),
	}
}

func UserFromEntities(user *entities.UserCreds) *User {
	if user == nil {
		return new(User)
	}
	return &User{
		Nickname: user.Nickname,
		Email:    user.Email,
		ID:       user.ID,
		Password: []byte(user.Password),
	}
}
