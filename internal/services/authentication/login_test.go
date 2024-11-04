package authentication

import (
	"context"
	"testing"
	"time"

	"github.com/blacksmith-vish/sso/internal/services/authentication/mocks"
	"github.com/blacksmith-vish/sso/internal/services/authentication/models"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	store_models "github.com/blacksmith-vish/sso/internal/store/models"
)

func Test_Login_Success(t *testing.T) {

	userID := uuid.New().String()
	appID := uuid.New().String()
	email := gofakeit.Email()
	password := randomPassword()
	passHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	request := models.LoginRequest{
		Email:    email,
		Password: password,
	}

	storeUser := store_models.User{
		ID:           userID,
		Email:        email,
		PasswordHash: passHash,
	}

	storeApp := store_models.App{
		Secret: appSecret,
		ID:     appID,
	}

	userProvider := mocks.NewUserProvider(t)
	userProvider.On(userProvider_UserByEmail, mock.Anything, email).Return(storeUser, nil)
	appProvider := mocks.NewAppProvider(t)
	appProvider.On(appProvider_App, mock.Anything, appID).Return(storeApp, nil)

	service := suite_NewService(nil, userProvider, appProvider)

	loginTime := time.Now()

	token, err := service.Login(context.TODO(), request, appID)
	require.NoError(t, err)

	tokenParsed, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, userID, claims["userID"].(string))
	assert.Equal(t, request.Email, claims["email"].(string))
	assert.Equal(t, appID, claims["appID"].(string))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(service.tokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)

}

func Test_Login_Fail_InvalidPassword(t *testing.T) {

	userID := uuid.New().String()
	appID := uuid.New().String()
	email := gofakeit.Email()
	password := randomPassword()
	passHash, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

	request := models.LoginRequest{
		Email:    email,
		Password: password,
	}

	storeUser := store_models.User{
		ID:           userID,
		Email:        email,
		PasswordHash: passHash,
	}

	storeApp := store_models.App{
		Secret: appSecret,
		ID:     appID,
	}

	userProvider := mocks.NewUserProvider(t)
	userProvider.On(userProvider_UserByEmail, mock.Anything, email).Return(storeUser, nil)
	appProvider := mocks.NewAppProvider(t)
	appProvider.On(appProvider_App, mock.Anything, appID).Return(storeApp, nil)

	service := suite_NewService(nil, userProvider, appProvider)

	token, err := service.Login(context.TODO(), request, appID)

	require.ErrorIs(t, err, models.ErrInvalidCredentials)
	require.Empty(t, token)
}

func Test_Login_Fail_App(t *testing.T) {

	noApp := store_models.App{}

	request := models.LoginRequest{
		Email:    gofakeit.Email(),
		Password: randomPassword(),
	}

	appID1 := uuid.New().String()
	appID2 := uuid.New().String()

	appProvider := mocks.NewAppProvider(t)
	appProvider.On(appProvider_App, mock.Anything, appID1).Return(noApp, store_models.ErrNotFound)
	appProvider.On(appProvider_App, mock.Anything, appID2).Return(noApp, errors.New("test error"))

	service := suite_NewService(nil, nil, appProvider)

	t.Run("invalid app ID", func(t *testing.T) {
		token, err := service.Login(context.TODO(), models.LoginRequest{}, WrongID)
		require.ErrorIs(t, err, models.ErrInvalidAppID)
		require.Empty(t, token)
	})

	t.Run("app not found / store returned apps.ErrAppNotFound", func(t *testing.T) {
		token, err := service.Login(context.TODO(), request, appID1)
		require.ErrorIs(t, err, models.ErrAppNotFound)
		require.Empty(t, token)
	})

	t.Run("app not found / other errors", func(t *testing.T) {
		token, err := service.Login(context.TODO(), request, appID2)
		require.ErrorIs(t, err, models.ErrAppsStore)
		require.Empty(t, token)
	})
}

func Test_Login_Fail_User(t *testing.T) {

	noUser := store_models.User{}

	appID := uuid.New().String()
	email1 := gofakeit.Email()
	email2 := gofakeit.Email()
	password := randomPassword()

	request1 := models.LoginRequest{
		Email:    email1,
		Password: password,
	}

	request2 := models.LoginRequest{
		Email:    email2,
		Password: password,
	}

	storeApp := store_models.App{
		Secret: appSecret,
		ID:     appID,
	}

	userProvider := mocks.NewUserProvider(t)
	userProvider.On(userProvider_UserByEmail, mock.Anything, email1).Return(noUser, store_models.ErrNotFound)
	userProvider.On(userProvider_UserByEmail, mock.Anything, email2).Return(noUser, errors.New("test error"))

	appProvider := mocks.NewAppProvider(t)
	appProvider.On(appProvider_App, mock.Anything, appID).Return(storeApp, nil)

	service := suite_NewService(nil, userProvider, appProvider)

	token, err := service.Login(context.TODO(), request1, appID)
	require.ErrorIs(t, err, models.ErrUserNotFound)
	require.Empty(t, token)

	token, err = service.Login(context.TODO(), request2, appID)
	require.ErrorIs(t, err, models.ErrUsersStore)
	require.Empty(t, token)

}

func Test_Login_Fail_InvalidUserData(t *testing.T) {

	appID := uuid.New().String()
	email := gofakeit.Email()
	password := randomPassword()

	service := suite_NewService(nil, nil, nil)

	TestingTable := []struct {
		name string
		// аргументы
		request models.LoginRequest
	}{
		{
			name: "empty password",
			request: models.LoginRequest{
				Email:    email,
				Password: "",
			},
		},
		{
			name: "empty email",
			request: models.LoginRequest{
				Email:    "",
				Password: password,
			},
		},
	}

	for _, tt := range TestingTable {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.Login(context.TODO(), tt.request, appID)
			require.ErrorIs(t, err, models.ErrInvalidRequest)
			require.Empty(t, token)
		})
	}
}
