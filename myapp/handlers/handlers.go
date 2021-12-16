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
