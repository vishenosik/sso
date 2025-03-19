package authentication

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/vishenosik/sso/internal/services/authentication/mocks"
	"github.com/vishenosik/sso/internal/services/authentication/models"
	store_models "github.com/vishenosik/sso/internal/store/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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
				err:     models.ErrUserNotFound,
			},
		},
		{
			name:   "not an admin, expect err",
			userID: ID4,
			expected: expectedResult{
				isAdmin: false,
				err:     models.ErrUsersStore,
			},
		},
		{
			name:   "invalid uuid",
			userID: WrongID,
			expected: expectedResult{
				isAdmin: false,
				err:     models.ErrUserInvalidID,
			},
		},
	}

	userProvider := mocks.NewUserProvider(t)

	userProvider.
		On(userProvider_IsAdmin, mock.Anything, ID1).Return(false, nil).
		On(userProvider_IsAdmin, mock.Anything, ID2).Return(true, nil).
		On(userProvider_IsAdmin, mock.Anything, ID3).Return(false, store_models.ErrNotFound).
		On(userProvider_IsAdmin, mock.Anything, ID4).Return(false, errors.New("test"))

	service := suite_NewService(nil, userProvider, nil)

	for _, tt := range TestingTable {
		t.Run(tt.name, func(t *testing.T) {
			isAdmin, err := service.IsAdmin(context.TODO(), tt.userID)
			assert.ErrorIs(t, err, tt.expected.err)
			assert.Equal(t, tt.expected.isAdmin, isAdmin)
		})
	}

}
