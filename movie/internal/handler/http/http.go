package http

import (
	"encoding/json"
	"errors"
	"log"
	"movieexample.com/movie/internal/controller"
	"net/http"
)

// Handler defines a movie handler.
type Handler struct {
	ctrl *controller.Controller
}

// New creates a new movie HTTP handler.
func New(ctrl *controller.Controller) *Handler {
	return &Handler{ctrl}
}

// GetMovieDetails handles GET /movie requests.
func (h *Handler) GetMovieDetails(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	details, err := h.ctrl.Get(req.Context(), id)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Repository get error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(details); err != nil {
		log.Printf("Response encode error: %v\n", err)
	}
}
