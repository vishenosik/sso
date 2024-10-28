package authentication

import (
	"context"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/blacksmith-vish/sso/internal/lib/config"
	"github.com/blacksmith-vish/sso/internal/services/authentication/mocks"
	auth_store "github.com/blacksmith-vish/sso/internal/store/sql/authentication"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func NewConfigTest() config.AuthenticationService {
	return config.AuthenticationService{
		TokenTTL: time.Minute,
	}
}

func Test_IsAdmin(t *testing.T) {

	type expectedResult struct {
		isAdmin bool
		err     error
	}

	ID1 := uuid.New().String()
	ID2 := uuid.New().String()
	ID3 := uuid.New().String()
	ID4 := "invalidx-uuid-xxxx-xxxx-xxxxxxxxxxxx"

	TestingTable := []struct {
		name string
		// аргументы
		userID string
		// ожидаемый результат
		expected expectedResult
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
				err:     auth_store.ErrUserNotFound,
			},
		},
		{
			name:   "invalid uuid",
			userID: ID4,
			expected: expectedResult{
				isAdmin: false,
				err:     ErrInvalidUserID,
			},
		},
	}

	// userSaver := mocks.NewUserSaver(t)
	userProvider := mocks.NewUserProvider(t)
	// appProvider := mocks.NewAppProvider(t)

	userProvider.
		On("IsAdmin", mock.Anything, ID1).
		Return(false, nil).
		On("IsAdmin", mock.Anything, ID2).
		Return(true, nil).
		On("IsAdmin", mock.Anything, ID3).
		Return(false, auth_store.ErrUserNotFound)

	service := NewService(
		slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelInfo})),
		NewConfigTest(),
		nil,
		userProvider,
		nil,
	)

	for _, tt := range TestingTable {
		t.Run(tt.name, func(t *testing.T) {
			isAdmin, err := service.IsAdmin(context.TODO(), tt.userID)
			assert.ErrorIs(t, err, tt.expected.err)
			assert.Equal(t, tt.expected.isAdmin, isAdmin)
		})
	}

}
