package controllers

import (
	"EnronEmailApi/services"
	"encoding/json"
	"fmt"
	"strconv"

	"net/http"

	"github.com/go-chi/chi/v5"
)

func IndexerEnron(w http.ResponseWriter, r *http.Request) {

	go services.IndexStart()

	w.Write([]byte("Cargando emails..."))
}

func SearchEmails(w http.ResponseWriter, r *http.Request) {
	text := chi.URLParam(r, "text")
	pageNum := r.URL.Query().Get("page")
	fmt.Println(pageNum)
	numero, err := strconv.Atoi(pageNum)

	fmt.Println(numero)
	if err != nil {
		http.Error(w, "Número de página inválido", http.StatusBadRequest)
		return
	}

	fmt.Println(text)

	respuesta := services.SearchEmails(&text, numero)

	json.NewEncoder(w).Encode(respuesta.Hits.Hits)
}
