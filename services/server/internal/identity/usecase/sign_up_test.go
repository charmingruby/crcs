package usecase_test

import (
	"context"
	"crcs/core/customerr"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/charmingruby/crcs/services/server/internal/identity/usecase"
)

func TestSignUp(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		input   usecase.SignUpInput
		name    string
		errType customerr.ErrorType
		wantErr bool
	}{
		{
			name: "success creates the user and returns a valid token",
			input: usecase.SignUpInput{
				FirstName: "Jane",
				LastName:  "Doe",
				Email:     "jane@example.com",
				Password:  "password123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestService(t)

			got, err := s.uc.SignUp(ctx, tt.input)

			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, customerr.IsType(err, tt.errType))
				assert.Empty(t, got.AccessToken)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, got.AccessToken)

			user, err := s.repo.FindByEmail(ctx, tt.input.Email)
			require.NoError(t, err)
			require.NotNil(t, user)
			assert.Equal(t, tt.input.FirstName, user.FirstName)
			assert.Equal(t, tt.input.LastName, user.LastName)
			assert.NotEqual(t, tt.input.Password, user.PasswordHash)
			require.NoError(t, bcrypt.CompareHashAndPassword(
				[]byte(user.PasswordHash),
				[]byte(tt.input.Password),
			))

			require.NoError(t, err)
		})
	}
}
