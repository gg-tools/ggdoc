package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	bindAddr = env("BIND_ADDR", ":80")
	htmlRoot = env("HTML_ROOT", "./statics/")
	docRoot  = filepath.Join(htmlRoot, "/swagger/")
	apiRoot  = filepath.Join(htmlRoot, "/apis/")
)

func main() {
	docFS := http.FileServer(http.Dir(docRoot))
	http.Handle("/docs/", http.StripPrefix("/docs/", docFS))

	apiFS := http.FileServer(http.Dir(apiRoot))
	http.Handle("/apis/", http.StripPrefix("/apis/", apiFS))

	indexPath := filepath.Join(htmlRoot, "index.html")
	t, err := template.ParseFiles(indexPath)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var apis []template.URL
		err := filepath.Walk(apiRoot, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			relativePath := strings.TrimPrefix(path, apiRoot)
			apiURI := filepath.ToSlash(filepath.Join("/apis/", relativePath))
			apis = append(apis, template.URL(apiURI))
			return nil
		})
		if err != nil {
			_, _ = w.Write([]byte(fmt.Sprintf("list apis failed: %s", err)))
			return
		}

		if err := t.Execute(w, apis); err != nil {
			_, _ = w.Write([]byte(fmt.Sprintf("list apis failed: %s", err)))
			return
		}
	})

	log.Fatal(http.ListenAndServe(bindAddr, nil))
}

func env(name string, defaultValue string) string {
	if v := os.Getenv(name); v != "" {
		return v
	} else {
		return defaultValue
	}
}
