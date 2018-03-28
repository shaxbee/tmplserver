package tmplserver

import (
	"io"
	"time"
)

// Resolver resolves resource into timestamp and reader
type Resolver interface {
	Resolve(name string) (time.Time, io.ReadSeeker, error)
}
