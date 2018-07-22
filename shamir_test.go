package shamir

import (
	"testing"
	"crypto/rand"

	"github.com/stretchr/testify/require"
)

func TestCreateCombine(t *testing.T) {
	require := require.New(t)
	// Short, medium, and long tests

	threshold := []int{10, 20, 100}
	shares := []int{21, 41, 201}

	for i := range threshold {
		secret, err := rand.Int(rand.Reader, P)
		require.NoError(err)

		created, err := CreateShares(threshold[i], shares[i], secret.Bytes())
		require.NoError(err)

		combined, err := CombineShares(created[:threshold[i]], threshold[i])
		require.NoError(err)

		require.Equal(combined, secret)
	}
}
