package internal

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/zishang520/socket.io/v2/socket"

	"github.com/willoma/keepakonf/frontend"
	"github.com/willoma/keepakonf/internal/client"
	"github.com/willoma/keepakonf/internal/data"
	"github.com/willoma/keepakonf/internal/log"
)

func Run(port int) (io.Closer, error) {
	io := socket.NewServer(nil, nil)
	logger := log.NewLogService(io.Sockets())
	data := data.New(io, logger)

	io.On("connection", func(clients ...any) {
		client.Serve(clients[0].(*socket.Socket), io.Sockets(), data, logger)
	})

	front, err := frontend.Handler()
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle("/socket.io/", io.ServeHandler(nil))
	mux.Handle("/", front)

	s := &http.Server{
		Addr:         fmt.Sprintf("127.0.0.1:%d", port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	slog.Info("Listening to TCP port", "port", port)

	go s.ListenAndServe()

	return s, nil
}
