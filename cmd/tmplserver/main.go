package main // import "github.com/shaxbee/tmplserver/cmd/tmplserver"

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"gopkg.in/yaml.v2"

	zglob "github.com/mattn/go-zglob"
	"github.com/pkg/errors"
	"github.com/shaxbee/tmplserver"
)

var (
	addr     = flag.String("addr", "0.0.0.0:8080", "Listen address")
	basePath = flag.String("base", ".", "Files base path")
	certFile = flag.String("cert", "", "TLS certificate path")
	keyFile  = flag.String("key", "", "TLS key path")
	dataFile = flag.String("data", "", "Template data file path (json/yaml)")
	env      = flag.Bool("env", false, "Load environment variables")
)

type tmplData struct {
	Env  map[string]interface{}
	Data map[string]interface{}
}

func parseEnv(src []string) map[string]interface{} {
	dst := make(map[string]interface{})

	for _, raw := range src {
		s := strings.SplitN(raw, "=", 2)
		k, v := s[0], s[1]

		if strings.Contains(v, ",") {
			dst[k] = strings.Split(v, ",")
		} else {
			dst[k] = v
		}
	}

	return dst
}

func load(base, data string, env bool) (*tmplData, error) {
	vals := tmplData{}
	if env {
		vals.Env = parseEnv(os.Environ())
	}

	if data == "" {
		return &vals, nil
	}

	ext := filepath.Ext(data)
	if ext != "json" && ext != "yaml" {
		return nil, errors.Errorf("Unsupported file type %s", ext)
	}

	b, err := ioutil.ReadFile(path.Join(base, data))
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to load %s", data)
	}

	dst := map[string]interface{}{}
	if err := yaml.Unmarshal(b, dst); err != nil {
		return nil, errors.Wrapf(err, "Failed to unmarshal %s", data)
	}

	vals.Data = dst
	return &vals, nil
}

func tmplResolver(base string, data string, env bool) (tmplserver.Resolver, error) {
	vals, err := load(base, data, env)
	if err != nil {
		return nil, err
	}

	files, err := zglob.Glob(path.Join(base, "**/*.tmpl"))
	if err != nil {
		return nil, errors.Wrap(err, "Invalid pattern")
	}

	if len(files) == 0 {
		return nil, nil
	}

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load templates")
	}

	return tmplserver.NewTmplResolver(tmpl, vals, time.Now()), nil
}

func main() {
	flag.Parse()

	base, err := filepath.Abs(*basePath)
	if err != nil {
		log.Fatalf("Invalid base path: %v", err)
	}

	resolvers := []tmplserver.Resolver{tmplserver.NewFileResolver(base)}
	tr, err := tmplResolver(base, *dataFile, *env)
	if err != nil {
		log.Print(err)
	}

	if tr != nil {
		resolvers = append(resolvers, tr)
	}

	l, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("Listening at %s", *addr)

	srv := tmplserver.NewServer(resolvers...)
	if *certFile != "" && *keyFile != "" {
		err = http.ServeTLS(l, srv, *certFile, *keyFile)
	} else {
		err = http.Serve(l, srv)
	}

	if err != http.ErrServerClosed && err != nil {
		log.Printf("Server error: %v", err)
	}
}
