package usecase_test

import (
	"context"
	"crcs/core/customerr"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/charmingruby/crcs/services/server/internal/auth/usecase"
)

func TestSignIn(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		input       usecase.SignInInput
		name        string
		errType     customerr.ErrorType
		wantErr     bool
		wantPersist bool
	}{
		{
			name: "unknown email returns unauthorized",
			input: usecase.SignInInput{
				Email:    "unknown@example.com",
				Password: "password123",
			},
			wantErr: true,
			errType: customerr.TypeUnauthorized,
		},
		{
			name: "wrong password returns unauthorized",
			input: usecase.SignInInput{
				Email:    "jane@example.com",
				Password: "wrong-password",
			},
			wantErr: true,
			errType: customerr.TypeUnauthorized,
		},
		{
			name: "success returns a valid token",
			input: usecase.SignInInput{
				Email:    "jane@example.com",
				Password: "password123",
			},
			wantPersist: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestService(t)

			if tt.wantPersist {
				signUp(t, s.uc)
			}

			got, err := s.uc.SignIn(ctx, tt.input)

			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, customerr.IsType(err, tt.errType))
				assert.Empty(t, got.AccessToken)
				return
			}

			require.NoError(t, err)
		})
	}
}
