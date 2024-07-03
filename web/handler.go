package web

import (
	"github.com/alekseibykov/goreddit"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"html/template"
	"net/http"
)

func NewHandler(store goreddit.Store) *Handler {
	h := &Handler{store: store, Mux: chi.NewMux()}

	h.Use(middleware.Logger)
	h.Route("/threads", func(r chi.Router) {
		r.Get("/", h.ThreadsList())
		r.Get("/new", h.ThreadsCreate())
		r.Post("/", h.ThreadsStore())
		r.Post("/{id}/delete", h.ThreadsDelete())
	})
	return h
}

type Handler struct {
	*chi.Mux
	store goreddit.Store
}

const threadsListHTML = `
<h1>Threads</h1>
<dl>
{{range .Threads}}
	<dt><strong>{{.Title}}</strong></dt>
	<dd>{{.Description}}</dd>
	<dd>
		<form action="/threads/{{.ID}}/delete" method="post">
			<button type="submit">Delete</button>
		</form>
	</dd>
{{end}}
</dl>
<a href="/threads/new">Create new thread</a>
`

func (h *Handler) ThreadsList() http.HandlerFunc {
	type data struct {
		Threads []goreddit.Thread
	}

	tmpl := template.Must(template.New("").Parse(threadsListHTML))
	return func(w http.ResponseWriter, r *http.Request) {
		tt, err := h.store.Threads()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, data{Threads: tt})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

const threadCreateHTML = `
<h1>New thread</h1>
<form action="/threads" method="post">
	<table>
		<tr>
			<td>Title</td>
			<td><input type="text" name="title" value="Test" /></td>
		</tr>
		<tr>
			<td>Description</td>
			<td><input type="text" name="description" value="Test" /></td>
		</tr>
	</table>
	<button type="submit">Create thread</button>
	<br>
	<a href="/threads/">Return to threads</a>
</form>
`

func (h *Handler) ThreadsCreate() http.HandlerFunc {
	tmpl := template.Must(template.New("").Parse(threadCreateHTML))
	return func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, nil)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *Handler) ThreadsStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("title")
		description := r.FormValue("description")

		if err := h.store.CreateThread(&goreddit.Thread{
			ID:          uuid.New(),
			Title:       title,
			Description: description,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/threads/", http.StatusFound)
	}
}

func (h *Handler) ThreadsDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = h.store.DeleteThread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/threads/", http.StatusFound)
	}
}
