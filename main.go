package main

import (
	"embed"
	"fmt"
	"github.com/gg-tools/apidoc-server/internal/markdown"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gg-tools/apidoc-server/internal/openapi"
)

//go:embed statics
var staticFS embed.FS

//go:embed pages
var pageFS embed.FS

var (
	bindAddr = env("BIND_ADDR", ":80")
	docsRoot = "docs/"
)

type DocEntry struct {
	App          string
	Title        string
	RelativePath template.URL
}

func main() {
	t, err := template.ParseFS(pageFS, "pages/*.html")
	if err != nil {
		log.Fatal(err)
	}
	statics := getStatics()

	http.Handle("/statics/", http.StripPrefix("/statics/", http.FileServer(statics)))
	http.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir(docsRoot))))
	http.HandleFunc("/mdview/", func(w http.ResponseWriter, r *http.Request) {
		filePath := strings.TrimPrefix(r.URL.Path, "/mdview")
		if err := t.ExecuteTemplate(w, "mdview.html", filePath); err != nil {
			_, _ = w.Write([]byte(fmt.Sprintf("rendor failed: %s", err)))
			return
		}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var entries []DocEntry
		err := filepath.Walk(docsRoot, func(path string, info fs.FileInfo, err error) error {
			switch {
			case info.IsDir():
				return nil
			case strings.HasSuffix(path, "_schema.yaml"):
				return nil
			case strings.HasSuffix(path, "_schema.json"):
				return nil
			}

			var title, app string
			switch {
			case strings.HasSuffix(path, ".json"), strings.HasSuffix(path, ".yaml"):
				app = "openapi"
				rd := openapi.NewDocumentReader(path)
				title, _ = rd.GetTitle()
			case strings.HasSuffix(path, ".md"):
				app = "markdown"
				rd := markdown.NewDocumentReader(path)
				title, _ = rd.GetTitle()
			default:
				app = "text"
				title = ""
			}

			relativePath := strings.TrimPrefix(path, docsRoot)
			apiURI := filepath.ToSlash(filepath.Join("/docs/", relativePath))
			if title == "" {
				title = relativePath
			}
			entry := DocEntry{
				App:          app,
				Title:        title,
				RelativePath: template.URL(apiURI),
			}
			entries = append(entries, entry)
			return nil
		})
		if err != nil {
			_, _ = w.Write([]byte(fmt.Sprintf("list entries failed: %s", err)))
			return
		}

		if err := t.ExecuteTemplate(w, "index.html", entries); err != nil {
			_, _ = w.Write([]byte(fmt.Sprintf("list entries failed: %s", err)))
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

func getStatics() http.FileSystem {
	var statics http.FileSystem

	stat, err := os.Stat("statics")
	if err == nil && stat.IsDir() {
		statics = http.FS(os.DirFS("statics"))
		log.Printf("using workdir statics")
	} else {
		embedStaticsSub, err := fs.Sub(staticFS, "statics")
		if err != nil {
			panic(err)
		}
		statics = http.FS(embedStaticsSub)
		log.Printf("[Main] using embeded statics")
	}

	return statics
}
