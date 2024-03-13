package store

import (
	"context"
	"math/rand"
	"os"
	"testing"

	"github.com/X-AROK/urlcut/internal/app/url"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TempFile() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	res := make([]byte, 10)
	for i := range res {
		res[i] = letters[rand.Intn(len(letters))]
	}

	return string(res)
}

func TestFileStore(t *testing.T) {
	dir := os.TempDir()
	fileStorePath := dir + string(os.PathSeparator) + TempFile()
	var fs *FileStore
	ctx := context.Background()
	t.Run("create file", func(t *testing.T) {
		store, err := NewFileStore(fileStorePath)
		require.NoError(t, err)
		assert.FileExists(t, fileStorePath)
		fs = store

		data, err := os.ReadFile(fileStorePath)
		require.NoError(t, err)
		assert.Empty(t, data)
	})
	t.Run("add to file", func(t *testing.T) {
		id, err := fs.Add(ctx, url.NewURL("https://practicum.yandex.ru"))
		require.NoError(t, err)
		assert.NotEmpty(t, id)

		data, err := os.ReadFile(fileStorePath)
		require.NoError(t, err)
		assert.Contains(t, string(data), id)
	})
	t.Run("close", func(t *testing.T) {
		err := fs.Close()
		require.NoError(t, err)
	})
	os.Remove(fileStorePath)
}
