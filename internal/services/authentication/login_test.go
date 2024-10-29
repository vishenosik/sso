package authentication

import (
	"context"
	"testing"
	"time"

	"github.com/blacksmith-vish/sso/internal/lib/jwt"
	"github.com/blacksmith-vish/sso/internal/services/authentication/mocks"
	"github.com/blacksmith-vish/sso/internal/services/authentication/models"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	store_models "github.com/blacksmith-vish/sso/internal/store/models"
)

const (
	emptyAppID = "iota"
	appID      = "a16fcc5e-d4de-4cf9-813f-e7ccf36f29d3"

	appSecret = "secret"

	passDefautlLen = 10
)

func Test_Login(t *testing.T) {

	// User(ctx context.Context, email string) (user models.User, err error)
	userProvider := mocks.NewUserProvider(t)
	// App(ctx context.Context, appID string) (app models.App, err error)
	appProvider := mocks.NewAppProvider(t)

	userID1 := uuid.New().String()
	email1 := gofakeit.Email()
	password1 := randomPassword()
	passHash1, _ := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)

	request1 := models.LoginRequest{
		Email:    email1,
		Password: password1,
	}

	if err := validator.New().Struct(request1); err != nil {
		panic(err)
	}

	storeUser1 := store_models.User{
		ID:           userID1,
		Email:        email1,
		PasswordHash: passHash1,
	}

	appID1 := uuid.New().String()

	storeApp1 := store_models.App{
		Secret: appSecret,
		ID:     appID1,
	}

	token1, _ := jwt.NewToken(storeUser1, storeApp1, time.Second)

	userProvider.On(userProvider_User, mock.Anything, email1).Return(storeUser1, nil)
	appProvider.On(appProvider_App, mock.Anything, appID1).Return(storeApp1, nil)

	// ID4 := "invalidx-uuid-xxxx-xxxx-xxxxxxxxxxxx"

	type expectedResult struct {
		token string
		err   error
	}

	TestingTable := []struct {
		name string
		// аргументы
		request models.LoginRequest
		appID   string
		// ожидаемый результат
		expected expectedResult
	}{
		{
			name:    "login success, no err",
			request: request1,
			appID:   appID1,

			expected: expectedResult{
				token: token1,
				err:   nil,
			},
		},
	}

	service := suite_NewService(nil, userProvider, appProvider)

	for _, tt := range TestingTable {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Login(context.TODO(), tt.request, tt.appID)
			assert.ErrorIs(t, err, tt.expected.err)
			// assert.Equal(t, tt.expected.token, token)
		})
	}

}

func randomPassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefautlLen)
}
