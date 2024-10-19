package main

import (
	"net/http"
	"path/filepath"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Static Files
	fileServer := http.FileServer(staticFileSystem{http.Dir(app.config.staticDir)})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamicMiddleware := alice.New(app.sessionManager.LoadAndSave)

	// Root
	mux.Handle("GET /{$}", dynamicMiddleware.ThenFunc(app.home))

	// Snippets Routes
	mux.Handle("GET /snippet/view/{id}", dynamicMiddleware.ThenFunc(app.getSnippetView))
	mux.Handle("GET /snippet/create", dynamicMiddleware.ThenFunc(app.getSnippetCreate))
	mux.Handle("POST /snippet/create", dynamicMiddleware.ThenFunc(app.postSnippetCreate))

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standardMiddleware.Then(mux)
}

type staticFileSystem struct {
	fs http.FileSystem
}

func (sfs staticFileSystem) Open(path string) (http.File, error) {
	file, err := sfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	if stat.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := sfs.fs.Open(index); err != nil {
			closeErr := file.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return file, nil
}
