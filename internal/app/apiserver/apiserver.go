package apiserver

import (
	"context"
	"github.com/dantedenis/reast-api-golang/internal/app/logger"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type APIServer struct {
	server    http.Server
	configure *config
	logger    *logger.Logger
	limiter   chan struct{}
}

func NewAPIServer(configure IBuild) *APIServer {
	con := configure.Build()

	api := &APIServer{
		configure: con,
		logger:    logger.NewLogger(os.Stdout, os.Stderr),
		limiter:   make(chan struct{}, 100),
	}
	api.server.Addr = con.GetAddr()

	return api
}

func (a *APIServer) Start(ctx context.Context) (err error) {
	a.configureRouter()

	connClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT)

		a.logger.InfoMsg.Printf("Server is shutting down to: %+v\n", <-sigint)
		if err = a.server.Shutdown(ctx); err != nil {
			a.logger.ErrorMsg.Println("Failed stop server, err:", err)
		}
		close(connClosed)
	}()

	go func() {
		if err = a.server.ListenAndServe(); err != http.ErrServerClosed {
			a.logger.ErrorMsg.Println("Failed start server")
		}
	}()

	a.logger.InfoMsg.Println("Server start to:", "http://"+a.server.Addr)

	<-connClosed

	a.logger.InfoMsg.Println("Server stop is successful")
	return
}

func (a *APIServer) configureRouter() {
	router := http.NewServeMux()

	router.HandleFunc("/", a.middleware(a.handlePost()))

	a.server.Handler = router
}

func (a *APIServer) handlePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			a.notFound(w)
			return
		}
		_, err := io.WriteString(w, "Hello")
		if err != nil {
			a.logger.ErrorMsg.Println("Error writer:", err)
		}
	}
}

func (a *APIServer) AddChan(dec int) {
	for i := 0; i < dec; i++ {
		a.limiter <- struct{}{}
	}
}

func (a *APIServer) DoneChan() {
	<-a.limiter
}
