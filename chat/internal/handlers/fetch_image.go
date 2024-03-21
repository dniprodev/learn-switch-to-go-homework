package handlers

import (
	"context"
	"io"
	"net/http"
)

func FetchImageHandler(fetchFunc func(context.Context) (io.Reader, error)) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		image, err := fetchFunc(r.Context())
		if err != nil {
			http.Error(w, "Failed to fetch image", http.StatusInternalServerError)
			return
		}

		// Write the image to the response.
		if _, err := io.Copy(w, image); err != nil {
			http.Error(w, "Failed to write image", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Image fetched successfully"))
	}
	return http.HandlerFunc(fn)
}
