package schema

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetMapEntry(t *testing.T) {
	t.Run("map[string]int", func(t *testing.T) {
		errs := []error{}
		var target int
		var unknown any = map[string]any{
			"value": 5,
		}

		GetMapEntry(unknown, "value", &target, &errs)

		err := errors.Join(errs...)
		require.NoError(t, err)
		require.Equal(t, 5, target)
	})

	t.Run("map[string]bool", func(t *testing.T) {
		errs := []error{}
		var target bool
		var unknown any = map[string]any{
			"value": true,
		}

		GetMapEntry(unknown, "value", &target, &errs)

		err := errors.Join(errs...)
		require.NoError(t, err)
		require.Equal(t, true, target)
	})

	t.Run("map[int]bool", func(t *testing.T) {
		errs := []error{}
		var target bool
		var unknown any = map[int]any{
			5: true,
		}

		GetMapEntry(unknown, 5, &target, &errs)

		err := errors.Join(errs...)
		require.NoError(t, err)
		require.Equal(t, true, target)
	})

	t.Run("not a map", func(t *testing.T) {
		errs := []error{}
		var target bool
		var unknown any = []string{}

		GetMapEntry(unknown, "value", &target, &errs)

		err := errors.Join(errs...)
		require.ErrorIs(t, err, ErrNotAMap)
	})

	t.Run("missing key", func(t *testing.T) {
		errs := []error{}
		var target bool
		var unknown any = map[string]any{
			"a": 5,
		}

		GetMapEntry(unknown, "b", &target, &errs)

		err := errors.Join(errs...)
		require.ErrorIs(t, err, ErrMissingKey)
	})

	t.Run("wrong value type", func(t *testing.T) {
		errs := []error{}
		var target bool
		var unknown any = map[string]any{
			"a": 5,
		}

		GetMapEntry(unknown, "a", &target, &errs)

		err := errors.Join(errs...)
		require.ErrorIs(t, err, ErrNotAMap)
	})
}
