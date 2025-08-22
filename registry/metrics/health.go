package metrics

import (
	"encoding/json"
	"net/http"
)

func Health(w http.ResponseWriter) {
	status := "ok"
	if status == "ok" {
		data := map[string]string{
			"status": "ok",
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(data)
	}
}
