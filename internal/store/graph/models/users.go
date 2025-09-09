package models

import "github.com/vishenosik/sso/internal/entities"

type User struct {
	Nickname     string   `json:"nickname,omitempty"`
	Email        string   `json:"email,omitempty"`
	ID           string   `json:"uuid,omitempty"`
	PasswordHash []byte   `json:"pass_hash,omitempty"`
	DType        []string `json:"dgraph.type,omitempty"`
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
		Password: string(user.PasswordHash),
	}
}

func UserFromEntities(user *entities.UserCreds) *User {
	if user == nil {
		return new(User)
	}
	return &User{
		Nickname:     user.Nickname,
		Email:        user.Email,
		ID:           user.ID,
		PasswordHash: []byte(user.Password),
	}
}
