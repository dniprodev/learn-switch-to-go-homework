package handlers

import (
	"context"
	"io"
	"net/http"
)

func ImageUploadHandler(uploadFunc func(context.Context, io.Reader) error) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		if r.Body == nil {
			http.Error(w, "Empty request body", http.StatusBadRequest)
			return
		}

		err := uploadFunc(r.Context(), r.Body)
		if err != nil {
			http.Error(w, "Failed to upload image", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Image uploaded successfully"))
	}
	return http.HandlerFunc(fn)
}
