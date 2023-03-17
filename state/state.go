package state

import (
	"context"
	"errors"
	"net/http"
)

// The following just fakes the sentry/datadog/db/logger pakages
// becaue I'm too lazy to come up with a more realistic example.
var (
	sentry  = sentryFake{}
	datadog = datadogFake{}
)

type logger struct{}
type datadogFake struct{}
type db struct{}
type sentryFake struct{}

func (l *logger) Error(msg interface{})
func (l *logger) Info(msg interface{})
func (l *logger) Shutdown()

// shows different kinds of cleanup functions
// these packages implement
func (dd *datadogFake) Shutdown() {}

func (db *db) Close() error {
	return nil
}

func (s *sentryFake) Flush() bool {
	return false
}

// end faking things

type Service struct {
	logger  *logger
	db      *db
	srv     *http.Server
	done    chan struct{}
	cleanup []func() error
}

func (s *Service) Shutdown(ctx context.Context) {
	<-ctx.Done()

	for _, cleaner := range s.cleanup {
		err := cleaner()
		s.logger.Error(err)
	}

	s.done <- struct{}{}
}

func (s *Service) Start() {
	// start your http server here
}

func (s *Service) Wait() {
	<-s.done
}

func NewService(logger *logger, db *db, srv *http.Server, otherShutdown ...func() error) Service {
	// NOTE: you'd probably pass these things in rather than
	// creating them here, I am again, mostly being lazy.
	cleanup := append(
		otherShutdown,
		func() error {
			// TODO: better context management
			return srv.Shutdown(context.TODO())
		},
		func() error {
			if flushed := sentry.Flush(); !flushed {
				// TODO static error
				return errors.New("unable to flush sentry at shutdown")
			}
			return nil
		},
		func() error {
			return db.Close()
		},
		// NOTE: always last so the logger is available
		// during shutdown.
		// Can also do this explicitly last in shutdown
		func() error {
			logger.Info("shutdown complete")
			logger.Shutdown()
			return nil
		},
	)

	// NOTE: register the signal handling here? or maybe in Start?

	return Service{
		logger:  logger,
		db:      db,
		srv:     srv,
		done:    make(chan struct{}),
		cleanup: cleanup,
	}
}

// elsewhere:
// The main thread becomes a lot clearer about what's going on
// and the communication between the components of shutdown are
// owned by the "Server" which is conceptually our application state
func main() {
	srv := NewService(&logger{}, &db{}, &http.Server{})

	// start the server and block the main go routing
	srv.Start()

	srv.Wait()
}
