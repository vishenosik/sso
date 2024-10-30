package authentication

import (
	"context"
	"testing"

	"github.com/blacksmith-vish/sso/internal/services/authentication/mocks"
	"github.com/blacksmith-vish/sso/internal/store/sql/components/users"
	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_IsAdmin(t *testing.T) {

	ID1 := uuid.New().String()
	ID2 := uuid.New().String()
	ID3 := uuid.New().String()
	ID4 := "invalidx-uuid-xxxx-xxxx-xxxxxxxxxxxx"

	type expectedResult struct {
		isAdmin bool
		err     error
	}

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
				err:     users.ErrUserNotFound,
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

	userProvider := mocks.NewUserProvider(t)

	userProvider.
		On(userProvider_IsAdmin, mock.Anything, ID1).Return(false, nil).
		On(userProvider_IsAdmin, mock.Anything, ID2).Return(true, nil).
		On(userProvider_IsAdmin, mock.Anything, ID3).Return(false, users.ErrUserNotFound)

	service := suite_NewService(nil, userProvider, nil)

	for _, tt := range TestingTable {
		t.Run(tt.name, func(t *testing.T) {
			isAdmin, err := service.IsAdmin(context.TODO(), tt.userID)
			assert.ErrorIs(t, err, tt.expected.err)
			assert.Equal(t, tt.expected.isAdmin, isAdmin)
		})
	}

}
