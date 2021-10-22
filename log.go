package errz

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/funcr"
)

const (
	logKeyCallstack = "callstack"
)

var logger logr.Logger = funcr.New(func(prefix, args string) {
	fmt.Println(prefix, args)
}, funcr.Options{})

// WithLogger set a package wide logger.
// Take a look at https://github.com/go-logr/logr
// for a list of popular logger implemetations.
func WithLogger(l logr.Logger) {
	logger = l
}
