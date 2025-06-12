package dto

import (
	"context"

	"github.com/vishenosik/sso-sdk/api"
	"github.com/vishenosik/sso/internal/entities"
	"github.com/vishenosik/sso/internal/services"
)

type AuthenticationDTO struct {
	service *services.AuthenticationService
}

func NewAuthenticationDTO(service *services.AuthenticationService) *AuthenticationDTO {
	return &AuthenticationDTO{
		service: service,
	}
}

func (s *AuthenticationDTO) LoginByEmail(
	ctx context.Context,
	email,
	password,
	appID string,
) (string, error) {
	return s.service.LoginByEmail(ctx, email, password, appID)
}

func (s *AuthenticationDTO) RegisterUser(
	ctx context.Context,
	user *api.User,
) (string, error) {
	return s.service.RegisterUser(ctx, userApiToEntities(user))
}

func (s *AuthenticationDTO) IsAdmin(
	ctx context.Context,
	userID string,
) (bool, error) {
	return s.service.IsAdmin(ctx, userID)
}

func userApiToEntities(user *api.User) *entities.UserCreds {
	return &entities.UserCreds{
		User: entities.User{
			ID:       user.ID,
			Nickname: user.Nickname,
			Email:    user.Email,
		},
		Password: user.Password,
	}
}
