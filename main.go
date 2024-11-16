package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/gg-tools/ggdoc/internal"
	"github.com/gg-tools/ggdoc/internal/markdown"
	"github.com/gg-tools/ggdoc/internal/openapi"
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
	App   string
	Title string
	URI   template.URL
}

func main() {
	t, err := template.ParseFS(pageFS, "pages/*.html")
	if err != nil {
		log.Fatal(err)
	}
	statics := getStatics()

	http.Handle("/statics/", http.StripPrefix("/statics/", http.FileServer(statics)))
	http.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir(docsRoot))))
	http.Handle("/mdview/", http.StripPrefix("/mdview/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if slices.Contains([]string{"/", "/index.html"}, r.URL.Path) || strings.HasSuffix(r.URL.Path, ".md") {
			filePath := filepath.Join("/", r.URL.Path)
			if err := t.ExecuteTemplate(w, "mdview.html", filePath); err != nil {
				_, _ = w.Write([]byte(fmt.Sprintf("rendor failed: %s", err)))
				return
			}
		} else {
			http.StripPrefix("docs", http.FileServer(http.Dir(docsRoot))).ServeHTTP(w, r)
		}
	})))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		entries, err := internal.DirectoryTree(docsRoot, func(path string, info fs.FileInfo) (*DocEntry, bool) {
			relPath, _ := filepath.Rel(docsRoot, path)
			relParts := strings.Split(relPath, string(filepath.Separator))
			name := relParts[len(relParts)-1]
			ext := filepath.Ext(name)

			var title, app string
			switch {
			case info.IsDir():
				return nil, true
			case strings.HasPrefix(name, "_"):
				return nil, false
			case slices.Contains([]string{".yaml", ".yml", ".json"}, ext):
				app = "openapi"
				rd := openapi.NewDocumentReader(path)
				title, _ = rd.GetTitle()
			case ext == ".md":
				app = "markdown"
				rd := markdown.NewDocumentReader(path)
				title, _ = rd.GetTitle()
			case ext == ".txt":
				app = "text"
				title = ""
			default:
				return nil, false
			}

			if title == "" {
				parts := strings.Split(path, string(filepath.Separator))
				title = parts[len(parts)-1]
			}

			uri := filepath.ToSlash(filepath.Join("/docs/", relPath))
			return &DocEntry{
				App:   app,
				Title: title,
				URI:   template.URL(uri),
			}, true
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
