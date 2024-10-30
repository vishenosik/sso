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
	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	store_models "github.com/blacksmith-vish/sso/internal/store/models"
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

func Test_Login_Fail_InvalidAppID(t *testing.T) {
	service := suite_NewService(nil, nil, nil)
	token, err := service.Login(context.TODO(), models.LoginRequest{}, WrongID)
	require.ErrorIs(t, err, ErrInvalidAppID)
	require.Empty(t, token)
}
func randomPassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefautlLen)
}
