package controllers

import (
	"EnronEmailApi/services"
	"encoding/json"
	"strconv"

	"net/http"

	"github.com/go-chi/chi/v5"
)

func IndexerEnron(w http.ResponseWriter, r *http.Request) {

	go services.IndexStart()

	w.Write([]byte("Loading emails..."))
}

func SearchEmails(w http.ResponseWriter, r *http.Request) {

	text := chi.URLParam(r, "text")
	pageNum := r.URL.Query().Get("page")

	numero, err := strconv.Atoi(pageNum)
	if err != nil {
		http.Error(w, "Invalid page number", http.StatusBadRequest)
		return
	}

	respuesta := services.SearchEmails(&text, numero)

	json.NewEncoder(w).Encode(respuesta.Hits.Hits)
}
