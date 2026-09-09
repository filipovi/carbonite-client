package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func newServer(r http.Handler, options options) *http.Server {
	var port int
	if options.port == nil {
		port = defaultHTTPPort
	} else {
		if *options.port == 0 {
			port = defaultHTTPPort + rand.IntN(1000)
		} else {
			port = *options.port
		}
	}

	var addr string
	if options.addr == nil || *options.addr == "" {
		addr = defaultAddr
	} else {
		addr = *options.addr
	}

	return &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("%s:%d", addr, port),
		IdleTimeout:  time.Minute,
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
	}
}

func (app *application) serve(opts ...Option) error {
	var options options
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			return err
		}
	}
	srv := newServer(app.routes(), options)

	// Create a shutdownError channel.
	// Used this to receive any errors returned by the graceful Shutdown() function.
	//
	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		log.Print("caught signal", map[string]string{
			"signal": s.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		log.Print("completing background tasks", map[string]string{
			"addr": srv.Addr,
		})

		// Call Wait() to block until our WaitGroup counter is zero
		app.wg.Wait()

		shutdownError <- nil
	}()

	log.Print("starting server", map[string]string{
		"addr": srv.Addr,
	})

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	log.Print("stopped server", map[string]string{
		"addr": srv.Addr,
	})
	return nil
}
