package db

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestNewInMemoryFromFile(t *testing.T) {
	t.Run("empty path", func(t *testing.T) {
		_, err := NewInMemoryFromFile("")

		assert.NoError(t, err)
	})
	t.Run("valid file", func(t *testing.T) {
		db, err := NewInMemoryFromFile("./testdata/sample_hosts.txt")

		assert.NoError(t, err)

		h1, err := db.Get("test")
		assert.NoError(t, err)
		assert.True(t, h1)
		h2, err := db.Get("test2")
		assert.NoError(t, err)
		assert.True(t, h2)
	})

	t.Run("invalid path", func(t *testing.T) {
		db, err := NewInMemoryFromFile("invalid")

		assert.ErrorIs(t, err, os.ErrNotExist)
		assert.Nil(t, db)
	})
}
