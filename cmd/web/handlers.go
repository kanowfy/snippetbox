package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/kanowfy/snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, err)
		return
	}

	data := &templateData{Snippets: snippets}
	app.render(w, r, "home.page.tmpl", data)

}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
	}

	data := &templateData{Snippet: snippet}
	app.render(w, r, "show.page.tmpl", data)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}



	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")
	// map to hold validation errors
	errors := make(map[string]string) 

	if strings.TrimSpace(title) == ""{
		errors["title"] = "this field can not be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		errors["title"] = "this field is too long (can not exceed 100 characters)"
	}

	if strings.TrimSpace(content) == "" {
		errors["content"] = "this field can not be blank"
	}

	if strings.TrimSpace(expires) == "" {
		errors["expires"] = "this field can not be blank"
	} else if expires != "365" && expires != "7" && expires != "1" {
		errors["expires"] = "this field is invalid"
	}

	// if any validation error, dump them in plain HTML response and return 
	if len(errors) > 0 {
		app.render(w, r, "create.page.tmpl", &templateData{
			FormData: r.PostForm,
			FormErrors: errors,
		})
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}

func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request){
	app.render(w, r, "create.page.tmpl", nil)
}