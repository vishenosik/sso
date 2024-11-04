package authentication

import (
	"context"
	"errors"
	"testing"

	"github.com/blacksmith-vish/sso/internal/services/authentication/mocks"
	"github.com/blacksmith-vish/sso/internal/services/authentication/models"
	store_models "github.com/blacksmith-vish/sso/internal/store/models"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_Register_Success(t *testing.T) {

	// userID := uuid.New().String()
	nickname := "nickname"
	email := gofakeit.Email()
	password := randomPassword()
	// passHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	request := models.RegisterRequest{
		Nickname: nickname,
		Email:    email,
		Password: password,
	}

	userSaver := mocks.NewUserSaver(t)
	userSaver.On(userSaver_SaveUser, mock.Anything, mock.Anything, nickname, email, mock.Anything).Return(nil)

	service := suite_NewService(userSaver, nil, nil)

	ID, err := service.RegisterNewUser(context.TODO(), request)
	require.NoError(t, err)

	err = validator.New().Var(ID, "required,uuid4")
	require.NoError(t, err)
}

func Test_Register_Fail_InvalidRequest(t *testing.T) {

	nickname := "nickname"
	email := gofakeit.Email()
	password := randomPassword()

	TestingTable := []struct {
		name string
		// аргументы
		request models.RegisterRequest
		// ожидаемый результат
		expectedErr error
	}{
		{
			name: "empty nickname",
			request: models.RegisterRequest{
				Nickname: "",
				Email:    email,
				Password: password,
			},
			expectedErr: models.ErrInvalidRequest,
		},
		{
			name: "empty email",
			request: models.RegisterRequest{
				Nickname: nickname,
				Email:    "",
				Password: password,
			},
			expectedErr: models.ErrInvalidRequest,
		},
		{
			name: "invalid email",
			request: models.RegisterRequest{
				Nickname: nickname,
				Email:    "email",
				Password: password,
			},
			expectedErr: models.ErrInvalidRequest,
		},
		{
			name: "empty password",
			request: models.RegisterRequest{
				Nickname: nickname,
				Email:    email,
				Password: "",
			},
			expectedErr: models.ErrInvalidRequest,
		},
		{
			name: "too long password (>72)",
			request: models.RegisterRequest{
				Nickname: nickname,
				Email:    email,
				Password: gofakeit.Password(true, true, true, true, false, 73),
			},
			expectedErr: models.ErrPasswordTooLong,
		},
	}

	service := suite_NewService(nil, nil, nil)

	for _, tt := range TestingTable {
		t.Run(tt.name, func(t *testing.T) {
			ID, err := service.RegisterNewUser(context.TODO(), tt.request)
			require.ErrorIs(t, err, tt.expectedErr)
			require.Empty(t, ID)
		})
	}
}

func Test_Register_Fail_Store(t *testing.T) {

	nickname1 := "nickname1"
	nickname2 := "nickname2"
	email := gofakeit.Email()
	password := randomPassword()

	request1 := models.RegisterRequest{
		Nickname: nickname1,
		Email:    email,
		Password: password,
	}

	request2 := models.RegisterRequest{
		Nickname: nickname2,
		Email:    email,
		Password: password,
	}

	userSaver := mocks.NewUserSaver(t)
	userSaver.On(userSaver_SaveUser, mock.Anything, mock.Anything, nickname1, email, mock.Anything).Return(store_models.ErrAlreadyExists)
	userSaver.On(userSaver_SaveUser, mock.Anything, mock.Anything, nickname2, email, mock.Anything).Return(errors.New("test"))

	service := suite_NewService(userSaver, nil, nil)

	ID, err := service.RegisterNewUser(context.TODO(), request1)
	require.ErrorIs(t, err, models.ErrUserExists)
	require.Empty(t, ID)

	ID, err = service.RegisterNewUser(context.TODO(), request2)
	require.ErrorIs(t, err, models.ErrUsersStore)
	require.Empty(t, ID)
}
