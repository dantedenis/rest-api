package apiserver

import (
	"context"
	"errors"
	conPkg "github.com/dantedenis/reast-api-golang/internal/app/config"
	"github.com/dantedenis/reast-api-golang/internal/app/logger"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	_ "time"
)

const (
	RateLimiter = 100
)

type APIServer struct {
	server    http.Server
	configure *conPkg.Config
	logger    logger.Logg
	limiter   int32
	mux       sync.Mutex
}

func NewAPIServer(configure conPkg.IBuild) *APIServer {
	con := configure.Build()

	api := &APIServer{
		configure: con,
		logger:    logger.NewLogger(os.Stdout, os.Stderr),
	}
	api.server.Addr = con.GetAddr()

	return api
}

func (a *APIServer) Start(ctx context.Context) (err error) {
	a.configureRouter()

	// Gorutine for
	connClosed := make(chan struct{})
	defer close(connClosed)
	go func() {
		sigint := make(chan os.Signal, 1)
		defer close(sigint)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT)

		a.logger.PrintfInfo("Server is shutting down to: %+v\n", <-sigint)
		a.mux.Lock()
		defer a.mux.Unlock()
		if err = a.server.Shutdown(ctx); err != nil {
			a.logger.PrintError("Failed stop server, err:", err.Error())
		}
		connClosed <- struct{}{}
	}()

	// Run server
	go func() {
		if err = a.server.ListenAndServe(); err != http.ErrServerClosed {
			a.logger.PrintError("Failed start server")
		}
	}()

	a.logger.PrintfInfo("Server start to:", "http://"+a.server.Addr)

	<-connClosed
	a.logger.PrintfInfo("Server stop is successful")
	return
}

func (a *APIServer) configureRouter() {
	router := http.NewServeMux()

	router.HandleFunc("/run", a.middleware(a.handler()))

	a.server.Handler = router
}

func (a *APIServer) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			a.clientError(w, http.StatusBadRequest, errors.New("invalid method"))
			return
		}

		if r.URL.Path != "/run" {
			a.notFound(w)
			return
		}

		if !a.AddQueue() {
			a.clientError(w, http.StatusTooManyRequests, errors.New("too many requestes"))
			return
		}
		defer a.DoneQueue()

		err := a.runner(w, r)
		if err != nil {
			a.clientError(w, http.StatusBadRequest, err)
			return
		}

	}
}

func (a *APIServer) AddQueue() bool {
	if atomic.LoadInt32(&a.limiter) < RateLimiter {
		atomic.AddInt32(&a.limiter, 1)
		return true
	}
	return false
}

func (a *APIServer) DoneQueue() {
	atomic.AddInt32(&a.limiter, -1)
}
