package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// The routes() method returns a servemux containing our application routes.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// dynamic: Unprotected routes
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))

	// Protected (authenticated-only) application routes, using a new "protected"
	// middleware chain which includes the requireAuthentication middleware.
	protected := dynamic.Append(app.requireAuthentication) // Authentification

	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))

	// Creates reusable middleware chain
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)
}
