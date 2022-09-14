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
	time.Sleep(time.Duration(delay) * time.Millisecond)

	if delay > 2000 {
		fmt.Printf("%s is raising an error!\n", name)
		select {
		case <-ctx.Done():
			// NOTE: don't try to send here, if the context is done,
			// then we closed the channel elsewhere and send would just panic
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

func finalizeFast(ctx context.Context, doneFuncs chan<- string, errors chan<- error) {
	simulateWork(ctx, "finalizeFast", 1, doneFuncs, errors)
}

func finalizeSlow(ctx context.Context, doneFuncs chan<- string, errors chan<- error) {
	simulateWork(ctx, "finalizeSlow", 1000, doneFuncs, errors)
}

func finalizeError(ctx context.Context, doneFuncs chan<- string, errors chan<- error) {
	simulateWork(ctx, "finalizeError", 5000, doneFuncs, errors)
}

func finalizeNever(ctx context.Context, doneFuncs chan<- string, errors chan<- error) {
	simulateWork(ctx, "finalizeError", 10*10*10*10*10, doneFuncs, errors)
}

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

	wg.Wait()

	// close the channel to indicate that we won't send any more
	select {
	case <-ctx.Done():
		// don't try to close these channels if the context is done, because we close
		// them elsewhere.
		break
	default:
		close(errors)
		close(doneFuncs)
	}

	// notify that we're done
	done <- true
	close(done)
}

func checkErrors(errors <-chan error) {
	// NOTE: we know in our example code that errors is closed
	// after all the finalizers run, so the for loop is safe here.
	// if errors is not closed though, this will loop indefinitely
	// so be careful with this construct.
	for err := range errors {
		fmt.Printf("    you an error: %s\n", err.Error())
	}
}

func checkCompleted(doneFuncs <-chan string) {
	// NOTE: ditto about this for loop!
	for name := range doneFuncs {
		fmt.Printf("    %s completed\n", name)
	}
}

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
			close(doneFuncs)
			close(errors)
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
