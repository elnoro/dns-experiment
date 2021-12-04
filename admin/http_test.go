package admin

import (
	"context"
	"github.com/elnoro/foxylock/m/v2/db"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestNewHttpServer(t *testing.T) {
	addr, tearDown := startTestServer()
	defer tearDown()

	r, err := http.Post("http://"+addr+"/addHost", "", strings.NewReader(`{"host":"test.test"}`))

	assert.NoError(t, err)
	assert.Equal(t, 201, r.StatusCode)

	body, err := io.ReadAll(r.Body)
	assert.NoError(t, err)
	assert.Equal(t, `{"host":"test.test"}`, string(body))

	r, err = http.Post("http://"+addr+"/removeHost", "", strings.NewReader(`{"host":"test.test"}`))
	assert.NoError(t, err)
	assert.Equal(t, 204, r.StatusCode)
}

func startTestServer() (string, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	address := "localhost:1234"
	srv := NewHttpServer(db.NewInMemory(), address)
	go func(ctx context.Context, srv DbServer) {
		_ = srv.Run(ctx)
	}(ctx, srv)

	time.Sleep(1 * time.Second)

	return address, cancel
}
