package services

import (

	// std
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	// pkg
	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	// internal pkg

	// internal

	"github.com/vishenosik/sso/internal/entities"
	store_models "github.com/vishenosik/sso/internal/store/models"
)

const (
	appSecret      = "secret"
	passDefautlLen = 10
	WrongID        = "invalidx-uuid-xxxx-xxxx-xxxxxxxxxxxx"

	userProvider_IsAdmin     = "IsAdmin"
	userProvider_UserByEmail = "UserByEmail"
	appProvider_App          = "App"
	userSaver_SaveUser       = "SaveUser"
)

func suite_NewService(
	t *testing.T,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
) *AuthenticationService {

	if userSaver == nil {
		userSaver = NewMockUserSaver(t)
	}

	if userProvider == nil {
		userProvider = NewMockUserProvider(t)
	}

	if appProvider == nil {
		appProvider = NewMockAppProvider(t)
	}

	service, err := NewAuthenticationService(
		userSaver,
		userProvider,
		appProvider,
		WithConfig(AuthenticationConfig{
			TokenTTL: time.Minute,
		}),
		WithLogger(slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo}))),
	)
	require.NoError(t, err)

	return service
}

func randomPassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefautlLen)
}

func Test_IsAdmin(t *testing.T) {

	t.Helper()
	t.Parallel()

	ID1 := uuid.New().String()
	ID2 := uuid.New().String()
	ID3 := uuid.New().String()
	ID4 := uuid.New().String()

	type expectedResult struct {
		err     error
		isAdmin bool
	}

	TestingTable := []struct {
		expected expectedResult
		name     string
		userID   string
	}{
		{
			name:   "not an admin, no err",
			userID: ID1,
			expected: expectedResult{
				isAdmin: false,
				err:     nil,
			},
		},
		{
			name:   "is admin, no err",
			userID: ID2,
			expected: expectedResult{
				isAdmin: true,
				err:     nil,
			},
		},
		{
			name:   "not an admin, expect err",
			userID: ID3,
			expected: expectedResult{
				isAdmin: false,
				err:     entities.ErrUserNotFound,
			},
		},
		{
			name:   "not an admin, expect err",
			userID: ID4,
			expected: expectedResult{
				isAdmin: false,
				err:     entities.ErrUsersStore,
			},
		},
		{
			name:   "invalid uuid",
			userID: WrongID,
			expected: expectedResult{
				isAdmin: false,
				err:     entities.ErrUserInvalidID,
			},
		},
	}

	userProvider := NewMockUserProvider(t)

	userProvider.
		On(userProvider_IsAdmin, mock.Anything, ID1).Return(false, nil).
		On(userProvider_IsAdmin, mock.Anything, ID2).Return(true, nil).
		On(userProvider_IsAdmin, mock.Anything, ID3).Return(false, store_models.ErrNotFound).
		On(userProvider_IsAdmin, mock.Anything, ID4).Return(false, errors.New("test"))

	service := suite_NewService(t, nil, userProvider, nil)

	for _, tt := range TestingTable {
		t.Run(tt.name, func(t *testing.T) {
			isAdmin, err := service.IsAdmin(context.TODO(), tt.userID)
			assert.ErrorIs(t, err, tt.expected.err)
			assert.Equal(t, tt.expected.isAdmin, isAdmin)
		})
	}

}

func Test_Login_Success(t *testing.T) {

	t.Helper()
	t.Parallel()

	userID := uuid.New().String()
	appID := uuid.New().String()
	email := gofakeit.Email()
	password := randomPassword()
	passHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	storeUser := &entities.UserCreds{
		User: entities.User{
			ID:    userID,
			Email: email,
		},
		PasswordHash: passHash,
	}

	storeApp := &entities.App{
		Secret: appSecret,
		ID:     appID,
	}

	userProvider := NewMockUserProvider(t)
	userProvider.On(userProvider_UserByEmail, mock.Anything, email).Return(storeUser, nil)

	appProvider := NewMockAppProvider(t)
	appProvider.On(appProvider_App, mock.Anything, appID).Return(storeApp, nil)

	service := suite_NewService(t, nil, userProvider, appProvider)

	loginTime := time.Now()

	token, err := service.LoginByEmail(context.TODO(), email, password, appID)
	require.NoError(t, err)

	tokenParsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, userID, claims["userID"].(string))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, claims["appID"].(string))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(service.config.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func Test_Login_Fail_InvalidPassword(t *testing.T) {

	t.Helper()
	t.Parallel()

	userID := uuid.New().String()
	appID := uuid.New().String()
	email := gofakeit.Email()
	password := randomPassword()
	passHash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	storeUser := &entities.UserCreds{
		User: entities.User{
			ID:    userID,
			Email: email,
		},
		PasswordHash: passHash,
	}

	storeApp := &entities.App{
		Secret: appSecret,
		ID:     appID,
	}

	userProvider := NewMockUserProvider(t)
	userProvider.On(userProvider_UserByEmail, mock.Anything, email).Return(storeUser, nil)

	appProvider := NewMockAppProvider(t)
	appProvider.On(appProvider_App, mock.Anything, appID).Return(storeApp, nil)

	service := suite_NewService(t, nil, userProvider, appProvider)

	token, err := service.LoginByEmail(context.TODO(), email, password, appID)

	require.ErrorIs(t, err, entities.ErrInvalidCredentials)
	require.Empty(t, token)
}

func Test_Login_Fail_App(t *testing.T) {

	t.Helper()
	t.Parallel()

	Email := gofakeit.Email()
	Password := randomPassword()

	appID1 := uuid.New().String()
	appID2 := uuid.New().String()

	appProvider := NewMockAppProvider(t)
	appProvider.On(appProvider_App, mock.Anything, appID1).Return(nil, store_models.ErrNotFound)
	appProvider.On(appProvider_App, mock.Anything, appID2).Return(nil, errors.New("test error"))

	service := suite_NewService(t, nil, nil, appProvider)

	t.Run("invalid app ID", func(t *testing.T) {
		token, err := service.LoginByEmail(context.TODO(), "", "", WrongID)
		require.ErrorIs(t, err, entities.ErrID)
		require.Empty(t, token)
	})

	t.Run("app not found / store returned apps.ErrAppNotFound", func(t *testing.T) {
		token, err := service.LoginByEmail(context.TODO(), Email, Password, appID1)
		require.ErrorIs(t, err, entities.ErrAppNotFound)
		require.Empty(t, token)
	})

	t.Run("app not found / other errors", func(t *testing.T) {
		token, err := service.LoginByEmail(context.TODO(), Email, Password, appID2)
		require.ErrorIs(t, err, entities.ErrAppsStore)
		require.Empty(t, token)
	})
}

func Test_Login_Fail_User(t *testing.T) {

	t.Helper()
	t.Parallel()

	appID := uuid.New().String()
	email1 := gofakeit.Email()
	email2 := gofakeit.Email()
	password := randomPassword()

	storeApp := &entities.App{
		Secret: appSecret,
		ID:     appID,
	}

	userProvider := NewMockUserProvider(t)
	userProvider.On(userProvider_UserByEmail, mock.Anything, email1).Return(nil, store_models.ErrNotFound)
	userProvider.On(userProvider_UserByEmail, mock.Anything, email2).Return(nil, errors.New("test error"))

	appProvider := NewMockAppProvider(t)
	appProvider.On(appProvider_App, mock.Anything, appID).Return(storeApp, nil)

	service := suite_NewService(t, nil, userProvider, appProvider)

	token, err := service.LoginByEmail(context.TODO(), email1, password, appID)
	require.ErrorIs(t, err, entities.ErrUserNotFound)
	require.Empty(t, token)

	token, err = service.LoginByEmail(context.TODO(), email2, password, appID)
	require.ErrorIs(t, err, entities.ErrUsersStore)
	require.Empty(t, token)

}

func Test_Login_Fail_InvalidUserData(t *testing.T) {

	t.Helper()
	t.Parallel()

	appID := uuid.New().String()
	email := gofakeit.Email()
	password := randomPassword()

	service := suite_NewService(t, nil, nil, nil)

	TestingTable := []struct {
		name string
		// аргументы
		email    string
		password string
		err      error
	}{
		{
			name:     "empty password",
			email:    email,
			password: "",
			err:      entities.ErrRequired,
		},
		{
			name:     "empty email",
			email:    "",
			password: password,
			err:      entities.ErrEmail,
		},
	}

	for _, tt := range TestingTable {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.LoginByEmail(context.TODO(), tt.email, tt.password, appID)
			require.ErrorIs(t, err, tt.err)
			require.Empty(t, token)
		})
	}
}

func Test_Register_Success(t *testing.T) {

	t.Helper()
	t.Parallel()

	nickname := "nickname"
	email := gofakeit.Email()
	password := randomPassword()

	user := &entities.UserCreds{
		User: entities.User{
			Nickname: nickname,
			Email:    email,
		},
		Password: password,
	}

	userSaver := NewMockUserSaver(t)
	userSaver.On(userSaver_SaveUser, mock.Anything, user).Return(nil)

	service := suite_NewService(t, userSaver, nil, nil)

	ID, err := service.RegisterUser(context.TODO(), user)
	require.NoError(t, err)

	err = validator.New().Var(ID, "required,uuid4")
	require.NoError(t, err)
}

func Test_Register_Fail_InvalidRequest(t *testing.T) {

	t.Helper()
	t.Parallel()

	nickname := "nickname"
	email := gofakeit.Email()
	password := randomPassword()

	TestingTable := []struct {
		expectedErr error
		request     entities.UserCreds
		name        string
	}{
		{
			name: "empty nickname",
			request: entities.UserCreds{
				User: entities.User{
					Nickname: "",
					Email:    email,
				},
				Password: password,
			},
			expectedErr: entities.ErrInvalidRequest,
		},
		{
			name: "empty email",
			request: entities.UserCreds{
				User: entities.User{
					Nickname: nickname,
					Email:    "",
				},
				Password: password,
			},
			expectedErr: entities.ErrInvalidRequest,
		},
		{
			name: "invalid email",
			request: entities.UserCreds{
				User: entities.User{
					Nickname: nickname,
					Email:    "email",
				},
				Password: password,
			},
			expectedErr: entities.ErrInvalidRequest,
		},
		{
			name: "empty password",
			request: entities.UserCreds{
				User: entities.User{
					Nickname: nickname,
					Email:    email,
				},
				Password: "",
			},
			expectedErr: entities.ErrInvalidRequest,
		},
		{
			name: "too long password (>72)",
			request: entities.UserCreds{
				User: entities.User{
					Nickname: nickname,
					Email:    email,
				},
				Password: gofakeit.Password(true, true, true, true, false, 73),
			},
			expectedErr: entities.ErrPasswordTooLong,
		},
	}

	service := suite_NewService(t, nil, nil, nil)

	for _, tt := range TestingTable {
		t.Run(tt.name, func(t *testing.T) {
			ID, err := service.RegisterUser(context.TODO(), &tt.request)
			require.ErrorIs(t, err, tt.expectedErr)
			require.Empty(t, ID)
		})
	}
}

func Test_Register_Fail_Store(t *testing.T) {

	t.Helper()
	t.Parallel()

	nickname1 := "nickname1"
	nickname2 := "nickname2"
	email := gofakeit.Email()
	password := randomPassword()

	request1 := entities.UserCreds{
		User: entities.User{
			Nickname: nickname1,
			Email:    email,
		},
		Password: password,
	}

	request2 := entities.UserCreds{
		User: entities.User{
			Nickname: nickname2,
			Email:    email,
		},
		Password: password,
	}

	userSaver := NewMockUserSaver(t)
	userSaver.On(userSaver_SaveUser, mock.Anything, &request1).Return(store_models.ErrAlreadyExists)
	userSaver.On(userSaver_SaveUser, mock.Anything, &request2).Return(errors.New("test"))

	service := suite_NewService(t, userSaver, nil, nil)

	ID, err := service.RegisterUser(context.TODO(), &request1)
	require.ErrorIs(t, err, entities.ErrUserExists)
	require.Empty(t, ID)

	ID, err = service.RegisterUser(context.TODO(), &request2)
	require.ErrorIs(t, err, entities.ErrUsersStore)
	require.Empty(t, ID)
}
