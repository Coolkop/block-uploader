package put

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	file, fh, err := r.FormFile("file")
	if err != nil {
		http.Error(w, fmt.Sprintf("get file from request"), http.StatusInternalServerError)

		return
	}

	data := make([]byte, fh.Size)
	_, err = file.Read(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("read file data"), http.StatusInternalServerError)

		return
	}

	uuid, err := h.processor.Process(r.Context(), data)
	if err != nil {
		http.Error(w, fmt.Sprintf("process upload data"), http.StatusInternalServerError)
		return
	}

	result := Result{ID: uuid.String()}
	responseData, err := json.Marshal(result)
	if err != nil {
		http.Error(w, fmt.Sprintf("prepare response data"), http.StatusInternalServerError)

		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(responseData)
}
