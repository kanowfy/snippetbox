package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// dynamic middleware so that it wont affect /static/ or any non dynamic routes
	dynamicMiddleware := alice.New(app.session.Enable, noSurf)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	// put create route first to avoid being mistaken as id
	mux.Post("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippet))
	mux.Get("/snippet/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createSnippetForm))
	mux.Get("/snippet/:id", dynamicMiddleware.ThenFunc(app.showSnippet))

	// user authentication
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.userSignupForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.userSignup))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.userLoginForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.userLogin))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.userLogout))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}
