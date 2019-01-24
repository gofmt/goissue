package api

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/render"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	workDir, _ := os.Getwd()
	staticDir := filepath.Join(workDir, "static")
	fileServer(r, "/static", http.Dir(staticDir))

	r.Route("/v1", func(r chi.Router) {
		r.Use(middleware.WithValue("api.version", "v1"))

		r.Get("/", wrap(func(w http.ResponseWriter, r *http.Request) error {
			return render.Render(w, r, ErrNotFound)
		}))

		r.Get("/test", wrap(func(w http.ResponseWriter, r *http.Request) error {
			return render.Render(w, r, ErrInvalidRequest(errors.New("invalid")))
		}))
	})

	return r
}

func wrap(h func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
		}
	}
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
