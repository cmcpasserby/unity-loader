package main

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/cmcpasserby/scli"
)

func createServeCmd() *scli.Command {
	return &scli.Command{
		Usage:         "serve [buildDirectory]",
		ShortHelp:     "Serves a WebGL Unity project",
		LongHelp:      "Serve a WebGL Unity project",
		ArgsValidator: scli.MaxArgs(1),
		Exec: func(ctx context.Context, args []string) error {
			var path string

			if len(args) > 0 {
				path = args[0]
			} else {
				wd, err := os.Getwd()
				if err != nil {
					return err
				}
				path = wd
			}

			fmt.Printf("Serving files in %q on port 8080\n", path)

			http.Handle("/", http.FileServer(http.Dir(path)))
			if err := http.ListenAndServe(":8080", makeGzipHandler(handler)); err != nil {
				return err
			}
			return nil
		},
	}
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func makeGzipHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			fn(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		gzr := gzipResponseWriter{Writer: gz, ResponseWriter: w}
		fn(gzr, r)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}
