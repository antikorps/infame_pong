package juego

import (
	_ "embed"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//go:embed html/partida.html
var paginaPartida string

//go:embed html/index.html
var index string

var mu sync.Mutex

func ServirIndex(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(index))
}

func ServirPaginaPartida(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	paginaPartidaPersonalizada := strings.Replace(paginaPartida, "##########", variables["identificador"], 1)

	w.WriteHeader(200)
	w.Write([]byte(paginaPartidaPersonalizada))
}

func gestionarComunicacionesWebSocket(manejadorWeb *ManejadorWeb, ws *websocket.Conn, identificadorPartida, identificadorUsuario string) {
	for {
		_, mensaje, mensajeError := ws.ReadMessage()
		if mensajeError != nil {
			// Un error se considera desconexion
			manejadorWeb.CerrarPartida(identificadorPartida, identificadorUsuario)
			return
		}
		if string(mensaje) == "iniciar" {
			manejadorWeb.IniciarJuego(identificadorPartida)
		}
		if string(mensaje) == "arriba" {
			manejadorWeb.MoverArriba(identificadorPartida, identificadorUsuario)
		}
		if string(mensaje) == "abajo" {
			manejadorWeb.MoverAbajo(identificadorPartida, identificadorUsuario)
		}
	}
}

func ManejarConexion(w http.ResponseWriter, r *http.Request, manejadorWeb *ManejadorWeb) {

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	identificadorPartida := r.URL.Query().Get("partida")
	identificadorUsuario := r.Header.Get("Sec-Websocket-Key")

	partidaExiste, partidaNumeroJugadores := manejadorWeb.ExistePartida(identificadorPartida)
	if partidaExiste {
		// Si existe una partida tiene que haber un jugador como mínimo
		// Si solo hay 1 se acepta y si hay 2 se rechaza la conexión
		if partidaNumeroJugadores > 1 {
			var enviarMensajeWs EnviarMensajeWS
			enviarMensajeWs.Mensaje = "partida completa"
			mu.Lock()
			ws.WriteJSON(&enviarMensajeWs)
			mu.Unlock()
			return
		}
		// Solo hay 1 jugador, añadir nuevo
		var jugador Jugador
		jugador.Identificador = identificadorUsuario
		jugador.Conexion = ws
		manejadorWeb.AñadirJugadorPartida(jugador, identificadorPartida)
		gestionarComunicacionesWebSocket(manejadorWeb, ws, identificadorPartida, identificadorUsuario)
		return
	}

	// No existe la partida, crear partida e insertar jugador
	var partida Partida
	partida.Identificador = identificadorPartida
	var jugador Jugador
	jugador.Identificador = identificadorUsuario
	jugador.Conexion = ws
	jugador.Principal = true

	partida.Jugadores = append(partida.Jugadores, jugador)
	manejadorWeb.AñadirPartida(partida)

	var enviarMensajeWs EnviarMensajeWS
	enviarMensajeWs.Mensaje = "primer jugador"
	mu.Lock()
	ws.WriteJSON(&enviarMensajeWs)
	mu.Unlock()

	gestionarComunicacionesWebSocket(manejadorWeb, ws, identificadorPartida, identificadorUsuario)

}
