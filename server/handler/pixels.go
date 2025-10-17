package handler

import (
	"encoding/json"
	"net/http"
	"tomashevich/server/database"
	"tomashevich/server/middleware"
)

func RegisterPixels(m *http.ServeMux, db *database.Database) {
	listPixels(m, db)
	paintPixel(m, db)
}

type listPixelsResponse struct {
	Pixels []database.Pixel `json:"pixels"`
}

func listPixels(m *http.ServeMux, db *database.Database) {
	m.HandleFunc("GET /pixels", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		pixels, err := db.GetPixels(r.Context())
		if err != nil {
			http.Error(w, "Can get pixels", http.StatusInternalServerError)
			return
		}

		if len(pixels) == 0 {
			pixels = make([]database.Pixel, 0)
		}

		json.NewEncoder(w).Encode(listPixelsResponse{pixels})
	})
}

type paintPixelData struct {
	X     int    `json:"x"`
	Y     int    `json:"y"`
	Color string `json:"color"`
}

func paintPixel(m *http.ServeMux, db *database.Database) {
	m.HandleFunc("POST /pixels:paint", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := middleware.GetSoulID(r.Context())
		if id == 0 {
			http.Error(w, "cant get your soul", http.StatusInternalServerError)
			return
		}

		var data paintPixelData
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "invalid form", http.StatusUnprocessableEntity)
			return
		}

		defer r.Body.Close()

		validColors := map[string]bool{
			"black":  true,
			"white":  true,
			"red":    true,
			"green":  true,
			"blue":   true,
			"yellow": true,
			"purple": true,
			"orange": true,
		}
		if _, ok := validColors[data.Color]; !ok {
			http.Error(w, "invalid color", http.StatusUnprocessableEntity)
			return
		}

		if data.X < 0 || data.Y < 0 {
			http.Error(w, "invalid x/y", http.StatusUnprocessableEntity)
			return
		}

		if err := db.PaintPixel(r.Context(), id, data.X, data.Y, data.Color); err != nil {
			http.Error(w, "cant paint this pixel", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
