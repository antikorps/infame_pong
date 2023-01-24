package juego

import (
	"sync"

	"github.com/gorilla/websocket"
)

type ManejadorWeb struct {
	sync.Mutex
	Partidas []Partida
}
type Partida struct {
	Identificador string      `json:"identificador"`
	Jugadores     []Jugador   `json:"jugadores"`
	Pelota        Pelota      `json:"pelota"`
	Canal         chan string `json:"-"`
}
type Jugador struct {
	Identificador string          `json:"id"`
	Principal     bool            `json:"principal"`
	Conexion      *websocket.Conn `json:"-"`
	X             int64           `json:"x"`
	Y             int64           `json:"y"`
	Ancho         int64           `json:"ancho"`
	Alto          int64           `json:"alto"`
	Puntos        int             `json:"puntos"`
}

type Pelota struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}

type EnviarMensajeWS struct {
	Mensaje string  `json:"mensaje"`
	Partida Partida `json:"partida"`
}
