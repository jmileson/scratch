package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	errFinalize    = errors.New("didn't finalize")
	finalizerFuncs = []finalizeFunc{
		finalizeFast,
		finalizeSlow,
		finalizeError,
		finalizeNever,
	}
)

type finalizeFunc func(context.Context, chan<- string, chan<- error)

func simulateWork(ctx context.Context, name string, delay int, doneFuncs chan<- string, errors chan<- error) {
	fmt.Printf("simulating work in %s\n", name)
	select {
	case <-ctx.Done():
		return
	case <-time.After(time.Duration(delay) * time.Millisecond):
		// This branch simulates the "work" finishing before the context is cancelled
	}

	// simulate some kind of error occurring during finalization
	if delay > 2000 {
		select {
		case <-ctx.Done():
			return
		default:
			errors <- errFinalize
		}
	}

	// NOTE: ditto here
	select {
	case <-ctx.Done():
		return
	default:
		doneFuncs <- name
	}
}

// finalizeFast completes quickly and doesn't return an error
func finalizeFast(ctx context.Context, doneFuncs chan<- string, errors chan<- error) {
	simulateWork(ctx, "finalizeFast", 1, doneFuncs, errors)
}

// finalizeSlow completes slowly, but doesn't return an error
func finalizeSlow(ctx context.Context, doneFuncs chan<- string, errors chan<- error) {
	simulateWork(ctx, "finalizeSlow", 1000, doneFuncs, errors)
}

// finalizeError completes slowly, and returns an error
func finalizeError(ctx context.Context, doneFuncs chan<- string, errors chan<- error) {
	simulateWork(ctx, "finalizeError", 5000, doneFuncs, errors)
}

// finalizeNever doesn't complete before the timeout, and would return an error if it could
// complete (but again, it can't because time).
func finalizeNever(ctx context.Context, doneFuncs chan<- string, errors chan<- error) {
	simulateWork(ctx, "finalizeError", 10*10*10*10*10, doneFuncs, errors)
}

func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	done := make(chan bool)
	go func() {
		defer close(done)
		wg.Wait()
	}()
	select {
	case <-done:
		// wg completed normally
		return false
	case <-time.After(timeout):
		return true
	}
}

// handleSig is responsible for handling signals sent from the OS and managing
// the running of finalization functions.
func handleSig(ctx context.Context, sig <-chan os.Signal, done chan<- bool, doneFuncs chan<- string, errors chan<- error) {
	// block until we receive a signal
	s := <-sig
	fmt.Printf("got you a signal: %s\n", s.String())

	var wg sync.WaitGroup
	for _, f := range finalizerFuncs {
		wg.Add(1)
		f := f
		go func() {
			defer wg.Done()
			f(ctx, doneFuncs, errors)
		}()
	}
	t, ok := ctx.Deadline()
	if !ok {
		// just choose something, this should never happen
		t = time.Now().Add(5 * time.Second)
	}
	timeout := time.Until(t)
	waitTimeout(&wg, timeout)

	// close the channel to indicate that we won't send any more
	close(errors)
	close(doneFuncs)

	// notify that we're done
	done <- true
	close(done)
}

// checkErrors prints any errors sent to the channel.
func checkErrors(errors <-chan error) {
	for err := range errors {
		// since we can't guarantee it closed, we check for zero values and treat them
		// like they're closed
		if err == nil {
			return
		}
		fmt.Printf("    you an error: %s\n", err.Error())

	}
}

// checkCompleted prints any completed finalizer functions sent to the channel.
func checkCompleted(doneFuncs <-chan string) {
	for name := range doneFuncs {
		if name == "" {
			return
		}
		fmt.Printf("    %s completed\n", name)

	}
}

// setupSignalHandling creates the necessary channels and registers signals to be handled.
func setupSignalHandling(ctx context.Context) (chan bool, chan string, chan error) {
	// handle CTRL+C interrupts
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	done := make(chan bool)
	// NOTE: these channels are buffered so that finalizers can send to them
	// without blocking. This is important because each finalizer must send
	// before finishing, and if more than 1 tries to send
	// on an unbuffered channel, any after the first will block until some
	// goroutine receives from it. Unfortunately, the function that receives
	// from the channels doesn't run until all the finalizers have finished.
	// So: if these channels are not buffered, we have a deadlock.
	doneFuncs := make(chan string, len(finalizerFuncs))
	errors := make(chan error, len(finalizerFuncs))

	// handle signals in the background
	// can't run on main goroutine
	go handleSig(ctx, sig, done, doneFuncs, errors)

	return done, doneFuncs, errors
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	done, doneFuncs, errors := setupSignalHandling(ctx)

	// block "forever", until we're done!
	for {
		fmt.Println("simulating work!")
		select {
		case <-ctx.Done():
			// if we get here before we can receive from done, then
			// we reached our timeout and graceful shutdown isn't an option.
			// So we should check and report which finalizers finished and
			// which didn't, and what errors we got for finished ones.
			// close(doneFuncs)
			// close(errors)
			fmt.Println("oh no, we timed out :(")
			fmt.Println("let's check for completed finalizers:")
			checkCompleted(doneFuncs)
			fmt.Println("here are the errors encountered:")
			checkErrors(errors)
			return
		case <-done:
			fmt.Println("we're done!")
			checkCompleted(doneFuncs)
			fmt.Println("but we should check for errors!")
			checkErrors(errors)
			return
		case <-time.After(1 * time.Second):
			// this is just a sleep to simulate work happening on the main thread,
			// e.g. an HTTP server or some other long running process.
		}
	}
}
