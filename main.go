package main

import (
	"log"
	"net/http"
	"time"

	_ "embed"

	"infame_pong/juego"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	manejadorWeb := juego.ManejadorWeb{}

	r.HandleFunc("/", juego.ServirIndex)
	r.HandleFunc("/partida/{identificador}", juego.ServirPaginaPartida)
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		juego.ManejarConexion(w, r, &manejadorWeb)
	})
	servidor := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(servidor.ListenAndServe())
}
