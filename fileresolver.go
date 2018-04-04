package tmplserver

import (
	"io"
	"os"
	"path"
	"time"
)

// NewFileResolver creates static files resolver
func NewFileResolver(base string) Resolver {
	return fileResolver{base}
}

type fileResolver struct {
	base string
}

func (fr fileResolver) Resolve(name string) (time.Time, io.ReadSeeker, error) {
	f, err := os.Open(path.Join(fr.base, name))

	switch {
	case os.IsNotExist(err):
		return time.Time{}, nil, nil
	case err != nil:
		return time.Time{}, nil, err
	}

	d, err := f.Stat()
	if err != nil {
		return time.Time{}, nil, err
	}

	return d.ModTime(), f, nil
}
