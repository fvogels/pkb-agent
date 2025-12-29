//go:build test

package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWords(t *testing.T) {
	t.Run("a", func(t *testing.T) {
		actual := Words("a")
		expected := []string{"a"}

		require.Equal(t, expected, actual)
	})

	t.Run("aaa", func(t *testing.T) {
		actual := Words("aaa")
		expected := []string{"aaa"}

		require.Equal(t, expected, actual)
	})

	t.Run("abc", func(t *testing.T) {
		actual := Words("abc")
		expected := []string{"abc"}

		require.Equal(t, expected, actual)
	})

	t.Run("a a", func(t *testing.T) {
		actual := Words("a a")
		expected := []string{"a", "a"}

		require.Equal(t, expected, actual)
	})

	t.Run("a abc", func(t *testing.T) {
		actual := Words("a abc")
		expected := []string{"a", "abc"}

		require.Equal(t, expected, actual)
	})

	t.Run(" a", func(t *testing.T) {
		actual := Words(" a")
		expected := []string{"a"}

		require.Equal(t, expected, actual)
	})

	t.Run("a ", func(t *testing.T) {
		actual := Words("a ")
		expected := []string{"a"}

		require.Equal(t, expected, actual)
	})

	t.Run(" a ", func(t *testing.T) {
		actual := Words(" a ")
		expected := []string{"a"}

		require.Equal(t, expected, actual)
	})

	t.Run("a  b", func(t *testing.T) {
		actual := Words("a  b")
		expected := []string{"a", "b"}

		require.Equal(t, expected, actual)
	})

	t.Run("", func(t *testing.T) {
		actual := Words("")
		expected := []string{}

		require.Equal(t, expected, actual)
	})

	t.Run("a123", func(t *testing.T) {
		actual := Words("a123")
		expected := []string{"a"}

		require.Equal(t, expected, actual)
	})

	t.Run("a,bc,def", func(t *testing.T) {
		actual := Words("a,bc,def")
		expected := []string{"a", "bc", "def"}

		require.Equal(t, expected, actual)
	})

	t.Run(" one two  three   four five", func(t *testing.T) {
		actual := Words(" one two  three   four five")
		expected := []string{"one", "two", "three", "four", "five"}

		require.Equal(t, expected, actual)
	})

	t.Run("ℛℛ ℛ", func(t *testing.T) {
		actual := Words("ℛℛ ℛ")
		expected := []string{"ℛℛ", "ℛ"}

		require.Equal(t, expected, actual)
	})
}
