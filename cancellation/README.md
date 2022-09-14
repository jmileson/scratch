# Cancellation

This example illustrates handling cancellation in Go functions when handling signals.

Run the package with

```golang
go run github.com/jmileson/scratch/cancellation
```

Then, while it's running, send a keyboard interrupt with `CTRL + C` to see the signal handling.

## Explanation

This is intended to illustrate how to shutdown a long running process like an HTTP server gracefully
when the OS sends a signal like SIGINT or SIGTERM the process. This is a common case to handle when
running application in K8s, since [K8s sends SIGTERM](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/)
to indicate that a pod will soon be forcefully shut down. In our example we're using SIGINT in place
of SIGTERM because it's easy to send that signal to a process with the keyboard.

The ideal flow looks like this:

- Run the binary
- Send SIGINT
- The process receives the SIGINT and runs all the finalizer functions
  - finalizer functions run in background because running them serially
    might prevent some of them from executing
  - once all the finalizer functions are done indicate to the main goroutine
    that finalization is done
- Once finalization is done
  - report the names of the functions that completed
  - report any errors encountered during finalization
- The process exits

However: in K8s, SIGTERM is followed by SIGKILL after some timeout - and SIGKILL can't be handled and
immediately terminates a process. We want to avoid loosing information about what finalizers ran and
errors that we got so we should try to ensure that the signal handling process completes PRIOR to this
hard timeout so we can report out before the process is killed dead. So the more real flow looks like:

- Run the binary
- Send SIGINT
- The process receives the SIGINT and runs all the finalizer functions
  - finalizer functions run in background because running them serially
    might prevent some of them from executing
  - once all the finalizer functions are done indicate to the main goroutine
    that finalization is done
- Once finalization is done OR our timeout is reached
  - report the names of the functions that completed
  - report any errors encountered during finalization
- The process exits
