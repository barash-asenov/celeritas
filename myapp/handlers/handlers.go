package handlers

import (
	"net/http"

	"github.com/barash-asenov/celeritas"
)

type Handlers struct {
	App *celeritas.Celeritas
}

func (h *Handlers) Home(rw http.ResponseWriter, r *http.Request) {
	err := h.App.Render.Page(rw, r, "home", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) GoPage(rw http.ResponseWriter, r *http.Request) {
	err := h.App.Render.GoPage(rw, r, "home", nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) JetPage(rw http.ResponseWriter, r *http.Request) {
	err := h.App.Render.JetPage(rw, r, "", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}
