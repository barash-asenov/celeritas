package handlers

import "net/http"

func (h *Handlers) UserLogin(rw http.ResponseWriter, r *http.Request) {
	err := h.App.Render.Page(rw, r, "login", nil, nil)

	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}
