package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/kanowfy/snippetbox/pkg/models"
)

type templateData struct {
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	CurrentYear int
}

// Custom template function
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// A FuncMap to hold the alias of the functions to pass in the template set
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// Template cache to early parse and avoid repeating parsing of templates
func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// Glob to get a slice of all page path
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extract the name (like home.page.tmpl)
		name := filepath.Base(page)
		// The template.FuncMap must be registered with the template set before
		// calling the ParseFiles() method. This means we have to use template.New() to
		// create an empty template set, use the Funcs() method to register the
		// template.FuncMap, and then parse the file as normal.

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// ParseGlob to add any layout to the template set
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add the parsed template set to cache
		cache[name] = ts
	}

	return cache, nil

}
