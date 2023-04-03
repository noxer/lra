package lra

import (
	"bytes"
	"io"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReaderAt(t *testing.T) {
	t.Run("is implemented", func(t *testing.T) {
		r := bytes.NewReader([]byte{1, 2, 3})
		rat := NewLazyReaderAt(r)

		require.Equal(t, r, rat)
	})

	t.Run("read successful", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(42))
		rat := NewLazyReaderAt(rnd)

		require.NotEqual(t, rnd, rat)

		buf1 := make([]byte, 32)
		n, err := rat.ReadAt(buf1, 100)
		require.EqualValues(t, 32, n)
		require.NoError(t, err)

		buf2 := make([]byte, 32)
		n, err = rat.ReadAt(buf2, 100)
		require.EqualValues(t, 32, n)
		require.NoError(t, err)

		require.ElementsMatch(t, buf1, buf2)

		n, err = rat.ReadAt(buf1, 20)
		require.EqualValues(t, 32, n)
		require.NoError(t, err)

		n, err = rat.ReadAt(buf2, 20)
		require.EqualValues(t, 32, n)
		require.NoError(t, err)

		require.ElementsMatch(t, buf1, buf2)
	})

	t.Run("read eof", func(t *testing.T) {
		rnd := rand.New(rand.NewSource(42))
		lim := io.LimitReader(rnd, 100)
		rat := NewLazyReaderAt(lim)

		buf1 := make([]byte, 32)
		n, err := rat.ReadAt(buf1, 200)
		require.Zero(t, n)
		require.Error(t, err)

		n, err = rat.ReadAt(buf1, 0)
		require.EqualValues(t, 32, n)
		require.NoError(t, err)

		n, err = rat.ReadAt(buf1, 99)
		require.EqualValues(t, 1, n)
		require.Error(t, err)
	})
}
