package get

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"file-storage/internal/services/get"
)

type Handler struct {
	processor processor
}

func New(processor processor) *Handler {
	return &Handler{processor: processor}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()

	fileID := chi.URLParam(r, "fileID")
	if fileID == "" {
		http.Error(w, "fileID in path required", http.StatusBadRequest)

		return
	}

	fileUUID, err := uuid.Parse(fileID)
	if err != nil {
		http.Error(w, "fileID has to be UUID format", http.StatusBadRequest)

		return
	}

	data, err := h.processor.Process(r.Context(), fileUUID)
	if err != nil {
		if errors.Is(err, get.NotFoundErr) {
			http.NotFound(w, r)
		} else {
			http.Error(w, fmt.Sprintf("process file obtaining"), http.StatusInternalServerError)
		}

		return
	}

	_, _ = w.Write(data)
}
