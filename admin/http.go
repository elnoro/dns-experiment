package admin

import (
	"context"
	"fmt"
	"github.com/elnoro/foxylock/m/v2/db"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type ginServer struct {
	db   db.HostDb
	addr string
}

func NewHttpServer(db db.HostDb, addr string) DbServer {
	return &ginServer{
		db:   db,
		addr: addr,
	}
}

type hostMessage struct {
	Host string `json:"host"`
}

func (h *ginServer) Run(ctx context.Context) error {
	router := gin.Default()
	router.POST("/addHost", func(c *gin.Context) {
		var m hostMessage
		err := c.BindJSON(&m)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})

			return
		}
		err = h.db.Save(m.Host)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(201, m)
	})

	router.POST("/removeHost", func(c *gin.Context) {
		var m hostMessage
		err := c.BindJSON(&m)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})

			return
		}

		err = h.db.Delete(m.Host)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})

			return
		}

		c.JSON(204, m)
	})

	srv := &http.Server{
		Addr:    h.addr,
		Handler: router,
	}

	errChan := make(chan error)
	go func(srv *http.Server, errChan chan<- error) {
		err := srv.ListenAndServe()
		errChan <- err
	}(srv, errChan)

	select {
	case err := <-errChan:
		return fmt.Errorf("gin server error, %w", err)
	case <-ctx.Done():
		return shutdown(srv)
	}
}

func shutdown(srv *http.Server) error {
	log.Print("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}
