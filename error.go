// Package errz provides highlevel error handling helpers
package errz

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Fatal panics on error and attaches a stack trace if needed.
// First parameter of msgs is used each following variadic arg is dropped
func Fatal(err error, msgs ...string) {
	if err != nil {
		var str string
		for _, msg := range msgs {
			str = msg
			break
		}

		if isPlain(err) {
			err = errors.Wrap(err, str)
		}
		panic(err)
	}
}

// Recover recovers a panic introduced by Fatal, any other function which calls panics
// or a memory corruption. Logs the error when called without args.
//
// Must be used at the top of the function defered
// defer Recover()
// or
// defer Recover(&err)
func Recover(errs ...*error) {
	var e *error
	for _, err := range errs {
		e = err
		break
	}

	if r := recover(); r != nil {
		err := r.(error)

		// In case of a memory corruption (e.g. nil pointer dereference)
		// the error does not have a stack trace and therefore one
		// is added here.
		if isPlain(err) {
			err = errors.WithStack(err)
		}

		// If error cant't bubble up -> log it
		if e != nil {
			*e = err
		} else {
			log(err)
		}
	}
}

// Log an error + stack trace to the underlying logger
// Usually used at the top level of a application to report errors without
// handling them.
func Log(err error, msgs ...string) {
	if err != nil {
		log(err, msgs...)
	}
}

// log splits the error in it's wrapped part and the original error message
// also writes it to the underlying logger.
func log(err error, msgs ...string) {
	var str string
	for _, msg := range msgs {
		str = msg
		break
	}
	cause := errors.Cause(err)
	sum := fmt.Sprintf("%+v", err)
	// remove original error from stack trace
	stack := strings.Replace(sum, cause.Error(), "", 1)

	// fmt.Printf("cause: %s\n", err.Error())
	// fmt.Printf("sum: %s\n", sum)
	// fmt.Printf("stack: %s\n", stack)
	logger.Error(cause, str, logKeyCallstack, stack)
}

// isPlain returns true when err is a plain error
// without a stacktrace.
func isPlain(err error) (plain bool) {

	cause := errors.Cause(err)
	sum := fmt.Sprintf("%+v", err)
	if cause.Error() == sum {
		plain = true
	}

	return plain
}
