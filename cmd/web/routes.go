package main

import (
	"net/http"
	"path/filepath"

	"github.com/jeremydwayne/snippets/ui"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Static Files
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	mux.Handle("GET /healthcheck", app.health())

	dynamicMiddleware := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// Root
	mux.Handle("GET /{$}", dynamicMiddleware.ThenFunc(app.home))

	// Snippets Routes
	mux.Handle("GET /snippets", dynamicMiddleware.ThenFunc(app.getSnippetIndex))
	mux.Handle("GET /snippet/view/{id}", dynamicMiddleware.ThenFunc(app.getSnippetView))

	// User Auth
	mux.Handle("GET /user/signup", dynamicMiddleware.ThenFunc(app.getUserSignup))
	mux.Handle("POST /user/signup", dynamicMiddleware.ThenFunc(app.postUserSignup))
	mux.Handle("GET /user/login", dynamicMiddleware.ThenFunc(app.getUserLogin))
	mux.Handle("POST /user/login", dynamicMiddleware.ThenFunc(app.postUserLogin))

	// Require Authentication
	protectedMiddleware := dynamicMiddleware.Append(app.requireAuthentication)
	mux.Handle("GET /snippet/create", protectedMiddleware.ThenFunc(app.getSnippetCreate))
	mux.Handle("POST /snippet/create", protectedMiddleware.ThenFunc(app.postSnippetCreate))
	mux.Handle("POST /user/logout", protectedMiddleware.ThenFunc(app.postUserLogout))

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
