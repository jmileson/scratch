/*
An example for using custom errors to handle bubbling up information
through your application without necessarily using concrete types
all over the place.

Advantages of this approach:
  - `error` (almost) everywhere
  - special handling isolated to creation and single function (or package)
  - custom error types can store whatever state needed
  - custom error types can implement any additional logic required

Disadvantages:
  - boilerplate
  - custom errors can quickly become complex, especially if you're wrapping other errors
  - HandleError can become a god function

Related reading:
  - https://pkg.go.dev/errors
  - https://go.dev/blog/error-handling-and-go
*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

// StatusAwareError knows what HTTP status code the
// error should return to a caller.
type StatusAwareError interface {
	// StatusAwareErrors must also be errors.
	error

	// Status returns the HTTP status code associated with the error.
	Status() int
}

// ErrorResponse indicates that an error occurred while
// processing a request.
type ErrorResponse struct {
	// HTTP Status code of the error.
	Code int

	// Message indicating what went wrong.
	Message string
}

// ValidationError indicates user provided
// inputs are invalid.
type ValidationError struct{}

// UniqueConstraintViolatedError indicates that
// a database save failed because a uniqe constraint
// on the table was violated.
type UniqueConstraintViolatedError struct{}

// ******************************
// implement our StatusAwareError
// ******************************

// Status implements StatusAwareError for ValidationError
func (e *ValidationError) Status() int {
	// 400 because the request was bad
	return http.StatusBadRequest
}

// Status implements StatusAwareError for ValidationError
func (e *UniqueConstraintViolatedError) Status() int {
	// 409 because the state of the system was violated
	return http.StatusConflict
}

// *****************************************
// implement error interface on custom types
// *****************************************

// Error implements error interface
// for ValidationError struct.
//
// See https://pkg.go.dev/builtin#error
// for the definition of this interface.
func (e *ValidationError) Error() string {
	return "the provided inputs are invalid for the following reasons: [foo, bar, baz]"
}

// Error implements error interface
// for UniqueConstraintViolatedError struct.
//
// See https://pkg.go.dev/builtin#error
// for the definition of this interface.
func (e *UniqueConstraintViolatedError) Error() string {
	return "cannot save record because another exists with the same ID"
}

// ************************************
// implement our error handler function
// ************************************

// HandleError handles any error raised by the application and
// creates an appropriate HTTP response.
func HandleError(w http.ResponseWriter, err error) {
	// 500 by default, as that represents an unknown state
	responseCode := http.StatusInternalServerError

	// if the error is status aware, then we can trust
	// it's status indication
	var statusAware StatusAwareError
	if errors.As(err, &statusAware) {
		responseCode = statusAware.Status()
	}

	// NOTE: or whatever - this is just an example of how to
	// somewhat generically handle errors that you want to
	// bubble up to users.
	resp := ErrorResponse{
		Code: responseCode,
		// NOTE: probably _not_ a good idea, this can leak a lot of info
		Message: err.Error(),
	}

	w.WriteHeader(responseCode)
	json.NewEncoder(w).Encode(&resp)
}

func main() {
	// validation error example:
	w := httptest.NewRecorder()
	HandleError(w, &ValidationError{})
	buf := strings.Builder{}
	io.Copy(&buf, w.Result().Body)
	fmt.Println(buf.String())

	// unique constraint violation example:
	w = httptest.NewRecorder()
	HandleError(w, &UniqueConstraintViolatedError{})
	buf = strings.Builder{}
	io.Copy(&buf, w.Result().Body)
	fmt.Println(buf.String())
}
