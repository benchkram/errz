// Package errz provides highlevel error handling helpers
package errz

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Fatal adds a stack trace to the error if needed and panics.
func Fatal(err error) {
	if err != nil {
		if isPlain(err) {
			err = errors.WithStack(err)
		}
		panic(err)
	}
}

// Fatalm adds a stack trace with an additonal message to the error and panics.
func Fatalm(err error, msg string) {
	if err != nil {
		err = errors.WithMessage(err, msgbegin)
		err = errors.WithMessage(err, msg)
		err = errors.WithMessage(err, msgend)
		if isPlain(err) {
			err = errors.WithStack(err)
		}
		panic(err)
	}
}

// FatalF adds a stack trace and an addtional message using printf semantics to the error and panics.
func Fatalf(err error, format string, args ...interface{}) {
	if err != nil {
		err = errors.WithMessage(err, msgbegin)
		err = errors.WithMessagef(err, format, args...)
		err = errors.WithMessage(err, msgend)
		if isPlain(err) {
			err = errors.WithStack(err)
		}
		panic(err)
	}
}

// Recover recovers a panic usually introduced by Fatal or a memory corruption.
// In the case of a plain error Recover decorates the error with a stacktrace.
// When a error pointer is passed `Recover(&err) the error is writted to the error pointer,
// and therefore returned to the calling function.
// In the case of `Recover()` called without args the error is logged to the underlying logr implementation.
//
// Example:
//
// func myfunc() (err error) {
//   defer errz.Recover(&err)
//   ...
//   ...
// }
//
func Recover(errs ...*error) {
	if r := recover(); r != nil {
		var e *error
		for _, err := range errs {
			e = err
			break
		}

		err := r.(error)

		// In case of a memory corruption (e.g. nil pointer dereference)
		// the error does not have a stack trace and therefore one
		// is added here.
		if isPlain(err) {
			err = errors.WithStack(err)
		}

		// If error cant't bubble up -> log it.
		if e != nil {
			*e = err
		} else {
			log(err)
		}
	}
}

// Log an error + stack trace to the underlying logger.
// Usually used at the top level of a application to log
// errors without handling them.
func Log(err error) {
	if err != nil {
		if isPlain(err) {
			err = errors.WithStack(err)
		}
		log(err)
	}
}

const msgbegin = "--annotationmsgbegin--"
const msgend = "--annotationmsgend--"

// log splits the error into the callstack, a message and the plain error
// then writes it to the underlying logger.
func log(err error) {

	var msg string
	var sanitizedStack string
	cause := errors.Cause(err)
	sum := fmt.Sprintf("%+v", err)

	// remove original error from stack trace
	stack := strings.Replace(sum, cause.Error(), "", 1)

	// The stack could have annotations attached. As github.com/pkg/errors does not allow
	// to access the msgs directly they are read from the summarized message.
	// Each message is wrapped by `msgbegin` and `msgend`
	scanner := bufio.NewScanner(bytes.NewBufferString(stack))
	var readingMsg bool
	for scanner.Scan() {
		txt := scanner.Text()
		if txt == msgbegin {
			readingMsg = true
			continue
		}
		if txt == msgend {
			readingMsg = false
			continue
		}
		if readingMsg {
			if msg != "" {
				msg = msg + "\n"
			}
			msg = msg + txt
			continue
		}
		sanitizedStack = sanitizedStack + txt + "\n"
	}

	logger.Error(cause, msg, logKeyCallstack, sanitizedStack)
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
