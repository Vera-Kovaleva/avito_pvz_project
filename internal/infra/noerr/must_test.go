package noerr_test

import (
	"errors"
	"testing"

	"avito_pvz/internal/infra/noerr"

	"github.com/stretchr/testify/require"
)

func Test(t *testing.T) {
	t.Parallel()

	err := errors.New("some err")

	require.PanicsWithError(t, err.Error(), func() {
		_ = noerr.Must(func() (int, error) {
			return 123, err
		}())
	})

	require.Equal(t, 123, noerr.Must(func() (int, error) {
		return 123, nil
	}()))
}
