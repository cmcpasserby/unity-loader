package serve

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

const (
	enableCors               = true
	enableWasmMultithreading = false
)

func Serve(path string, port int) error {
	fileServer := http.FileServer(http.Dir(path))
	handler := webGlHandler(fileServer)
	addr := fmt.Sprintf(":%d", port)

	log.Printf("Serving %s on http://localhost%s\n", path, addr)
	return http.ListenAndServe(addr, handler)
}

func webGlHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if enableCors {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}

		encodingExt := filepath.Ext(r.URL.Path)

		switch encodingExt {
		case ".br":
			w.Header().Set("Content-Encoding", "br")
		case ".gz":
			w.Header().Set("Content-Encoding", "gzip")
		}

		ext := filepath.Ext(strings.TrimSuffix(r.URL.Path, encodingExt))

		if enableWasmMultithreading && (r.URL.Path == "/" || ext == ".js" || ext == ".html" || ext == ".htm") {
			w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
			w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
			w.Header().Set("Cross-Origin-Resource-Policy", "cross-origin")
		}

		switch ext {
		case ".wasm":
			w.Header().Set("Content-Type", "application/wasm")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case ".json":
			w.Header().Set("Content-Type", "application/json")
		case ".bundle":
			fallthrough
		case ".unityweb":
			fallthrough
		case ".data":
			w.Header().Set("Content-Type", "application/octet-stream")
		}

		_, hasModifiedSince := r.Header["If-Modified-Since"]
		_, hasIfNoneMatch := r.Header["If-None-Match"]

		if r.Header.Get("Cache-Control") == "no-cache" && (hasModifiedSince || hasIfNoneMatch) {
			r.Header.Del("Cache-Control")
		}

		next.ServeHTTP(w, r)
	})
}
