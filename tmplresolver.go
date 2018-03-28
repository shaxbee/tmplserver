package tmplserver

import (
	"bytes"
	"io"
	"sync"
	"text/template"
	"time"
)

// NewTmplResolver creates template resolver
func NewTmplResolver(tmpl *template.Template, data interface{}, ts time.Time) Resolver {
	return &tmplResolver{
		tmpl:     tmpl,
		data:     data,
		ts:       ts,
		resolved: make(map[string]entry),
	}
}

type tmplResolver struct {
	sync.RWMutex
	tmpl     *template.Template
	data     interface{}
	ts       time.Time
	resolved map[string]entry
}

type entry struct {
	b   []byte
	err error
}

func (tr *tmplResolver) Resolve(name string) (time.Time, io.ReadSeeker, error) {
	name += ".tmpl"
	if tr.tmpl.Lookup(name) == nil {
		return time.Time{}, nil, nil
	}

	e := tr.lookup(name)
	if e.b == nil && e.err == nil {
		e = tr.execute(name)
	}

	switch {
	case e.err != nil:
		return time.Time{}, nil, e.err
	case e.b != nil:
		return tr.ts, bytes.NewReader(e.b), nil
	default:
		return time.Time{}, nil, nil
	}
}

func (tr *tmplResolver) lookup(name string) entry {
	tr.RLock()
	e := tr.resolved[name]
	tr.RUnlock()

	return e
}

func (tr *tmplResolver) execute(name string) entry {
	e := entry{}
	wr := &bytes.Buffer{}
	if err := tr.tmpl.ExecuteTemplate(wr, name, tr.data); err != nil {
		e.err = err
	} else {
		e.b = wr.Bytes()
	}

	tr.Lock()
	tr.resolved[name] = e
	tr.Unlock()

	return e
}
