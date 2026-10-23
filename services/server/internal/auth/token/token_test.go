package token_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/charmingruby/crcs/services/server/internal/auth/token"
)

func TestManagerSignAndVerify(t *testing.T) {
	manager := token.New("secret", time.Hour)

	signed, err := manager.Sign("user-1")
	require.NoError(t, err)
	require.NotEmpty(t, signed)

	userID, err := manager.Verify(signed)
	require.NoError(t, err)
	assert.Equal(t, "user-1", userID)
}

func TestManagerVerifyErrors(t *testing.T) {
	manager := token.New("secret", time.Hour)

	signed, err := manager.Sign("user-1")
	require.NoError(t, err)

	expiredManager := token.New("secret", -time.Minute)
	expired, err := expiredManager.Sign("user-1")
	require.NoError(t, err)

	otherManager := token.New("other-secret", time.Hour)

	tests := []struct {
		input string
		name  string
	}{
		{name: "malformed token", input: "not-a-token"},
		{name: "tampered token", input: signed + "tampered"},
		{name: "token signed with another secret", input: otherManagerSigned(t, otherManager)},
		{name: "expired token", input: expired},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := manager.Verify(tt.input)

			require.ErrorIs(t, err, token.ErrInvalidToken)
			assert.Empty(t, userID)
		})
	}
}

func TestManagerSignRequiresSubject(t *testing.T) {
	manager := token.New("secret", time.Hour)

	signed, err := manager.Sign("")
	require.NoError(t, err)

	userID, err := manager.Verify(signed)

	require.ErrorIs(t, err, token.ErrInvalidToken)
	assert.Empty(t, userID)
}

func otherManagerSigned(t *testing.T, manager *token.Manager) string {
	t.Helper()

	signed, err := manager.Sign("user-1")
	require.NoError(t, err)

	return signed
}
