package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/FabianAlmos/contfuncs-sdk/fn_http"
)

func Handle[In any, Out any](fn Handler[In, Out]) error {
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Handling request...")

		var req fn_http.Request[In]
		if r.ContentLength > 0 {
			if err := json.NewDecoder(r.Body).Decode(&req.Payload); err != nil {
				writeError(w, http.StatusBadRequest, fmt.Sprintf("failed to decode payload, err: %s", err))
				return
			}
		}

		resp, err := fn(r.Context(), req)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})

	log.Printf("Function listening on port: %q\n", port)
	return http.ListenAndServe(":"+port, nil)
}
