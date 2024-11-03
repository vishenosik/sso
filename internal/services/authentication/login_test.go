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
	"github.com/blacksmith-vish/sso/internal/store/sql/components/apps"
	"github.com/blacksmith-vish/sso/internal/store/sql/components/users"
)

const (
	appSecret      = "secret"
	passDefautlLen = 10
	WrongID        = "invalidx-uuid-xxxx-xxxx-xxxxxxxxxxxx"
)

func Test_Login_Success_TokenValid_NoErr(t *testing.T) {

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
	userProvider.On(userProvider_User, mock.Anything, email).Return(storeUser, nil)
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
	userProvider.On(userProvider_User, mock.Anything, email).Return(storeUser, nil)
	appProvider := mocks.NewAppProvider(t)
	appProvider.On(appProvider_App, mock.Anything, appID).Return(storeApp, nil)

	service := suite_NewService(nil, userProvider, appProvider)

	token, err := service.Login(context.TODO(), request, appID)

	require.ErrorIs(t, err, ErrInvalidCredentials)
	require.Empty(t, token)
}

func Test_Login_Fail_App(t *testing.T) {

	noApp := store_models.App{}
	noRequest := models.LoginRequest{}

	TestingErr := errors.New("test error")

	appID1 := uuid.New().String()
	appID2 := uuid.New().String()

	appProvider := mocks.NewAppProvider(t)
	appProvider.On(appProvider_App, mock.Anything, appID1).Return(noApp, apps.ErrAppNotFound)
	appProvider.On(appProvider_App, mock.Anything, appID2).Return(noApp, TestingErr)

	service := suite_NewService(nil, nil, appProvider)

	t.Run("invalid app ID", func(t *testing.T) {
		token, err := service.Login(context.TODO(), models.LoginRequest{}, WrongID)
		require.ErrorIs(t, err, ErrInvalidAppID)
		require.Empty(t, token)
	})

	t.Run("app not found / store returned apps.ErrAppNotFound", func(t *testing.T) {
		token, err := service.Login(context.TODO(), noRequest, appID1)
		require.ErrorIs(t, err, ErrInvalidAppID)
		require.Empty(t, token)
	})

	t.Run("app not found / other errors", func(t *testing.T) {
		token, err := service.Login(context.TODO(), noRequest, appID2)
		require.ErrorIs(t, err, TestingErr)
		require.Empty(t, token)
	})

}

func Test_Login_Fail_User(t *testing.T) {

	noUser := store_models.User{}
	TestingErr := errors.New("test error")

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
	userProvider.On(userProvider_User, mock.Anything, email1).Return(noUser, users.ErrUserNotFound)
	userProvider.On(userProvider_User, mock.Anything, email2).Return(noUser, TestingErr)

	appProvider := mocks.NewAppProvider(t)
	appProvider.On(appProvider_App, mock.Anything, appID).Return(storeApp, nil)

	service := suite_NewService(nil, userProvider, appProvider)

	token, err := service.Login(context.TODO(), request1, appID)
	require.ErrorIs(t, err, ErrInvalidCredentials)
	require.Empty(t, token)

	token, err = service.Login(context.TODO(), request2, appID)
	require.ErrorIs(t, err, TestingErr)
	require.Empty(t, token)

}

func Test_Login_Fail_InvalidUserData(t *testing.T) {

	appID := uuid.New().String()
	email := gofakeit.Email()
	password := randomPassword()

	storeApp := store_models.App{
		Secret: appSecret,
		ID:     appID,
	}

	appProvider := mocks.NewAppProvider(t)
	appProvider.On(appProvider_App, mock.Anything, appID).Return(storeApp, nil)

	service := suite_NewService(nil, nil, appProvider)

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
		{
			name: "invalid email",
			request: models.LoginRequest{
				Email:    "",
				Password: password,
			},
		},
	}

	for _, tt := range TestingTable {
		t.Run(tt.name, func(t *testing.T) {
			token, err := service.Login(
				context.TODO(),
				tt.request,
				appID,
			)

			require.ErrorIs(t, err, ErrInvalidCredentials)
			require.Empty(t, token)
		})
	}
}

func randomPassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefautlLen)
}
