package juego

import (
	"fmt"
	"math/rand"
	"time"
)

func (mw *ManejadorWeb) MoverPelota(identificadorPartida string) {

	desplazamiento := int64(5)

	opcionesX := []string{"izquierda", "derecha"}
	opcionesY := []string{"arriba", "abajo"}

	tendenciaX := opcionesX[rand.Intn(len(opcionesX))]
	tendenciaY := opcionesY[rand.Intn(len(opcionesY))]

	for {
		for i := range mw.Partidas {
			if mw.Partidas[i].Identificador == identificadorPartida {
				for j := range mw.Partidas[i].Jugadores {

					// Colisiones tamaño: 20x60
					// Jugador 1: cambiar tendencia derecha
					if mw.Partidas[i].Pelota.X <= mw.Partidas[i].Jugadores[0].X+20 {
						// Coincide en el eje X
						if mw.Partidas[i].Pelota.Y <= mw.Partidas[i].Jugadores[0].Y+60 {
							// Coincide en el eje Y
							tendenciaX = "derecha"
						}
					}
					// Jugador 2: cambiar tendencia derecha
					if mw.Partidas[i].Pelota.X >= mw.Partidas[i].Jugadores[1].X {
						// Coincide en el eje X
						if mw.Partidas[i].Pelota.Y <= mw.Partidas[i].Jugadores[1].Y+60 {
							// Coincide en el eje Y
							tendenciaX = "izquierda"
						}
					}

					if mw.Partidas[i].Pelota.X <= 5 {
						tendenciaX = "derecha"
						tendenciaY = opcionesY[rand.Intn(len(opcionesY))]
						mw.Lock()
						mw.Partidas[i].Jugadores[1].Puntos += 1
						mw.Unlock()
					}
					if mw.Partidas[i].Pelota.X >= 595 {
						tendenciaX = "izquierda"
						tendenciaY = opcionesY[rand.Intn(len(opcionesY))]
						mw.Lock()
						mw.Partidas[i].Jugadores[0].Puntos += 1
						mw.Unlock()
					}

					if mw.Partidas[i].Pelota.Y <= 5 {
						tendenciaY = "abajo"
						tendenciaX = opcionesX[rand.Intn(len(opcionesX))]
					}

					if mw.Partidas[i].Pelota.Y >= 375 {
						tendenciaY = "arriba"
						tendenciaX = opcionesX[rand.Intn(len(opcionesX))]
					}

					if tendenciaX == "izquierda" && tendenciaY == "arriba" {
						mw.Lock()
						mw.Partidas[i].Pelota.X -= desplazamiento
						mw.Partidas[i].Pelota.Y -= desplazamiento
						mw.Unlock()
					}
					if tendenciaX == "izquierda" && tendenciaY == "abajo" {
						mw.Lock()
						mw.Partidas[i].Pelota.X -= desplazamiento
						mw.Partidas[i].Pelota.Y += desplazamiento
						mw.Unlock()
					}
					if tendenciaX == "derecha" && tendenciaY == "arriba" {
						mw.Lock()
						mw.Partidas[i].Pelota.X += desplazamiento
						mw.Partidas[i].Pelota.Y -= desplazamiento
						mw.Unlock()
					}
					if tendenciaX == "derecha" && tendenciaY == "abajo" {
						mw.Lock()
						mw.Partidas[i].Pelota.X += desplazamiento
						mw.Partidas[i].Pelota.Y += desplazamiento
						mw.Unlock()
					}

					var enviarMensajeWs EnviarMensajeWS
					enviarMensajeWs.Mensaje = "jugar"
					enviarMensajeWs.Partida = mw.Partidas[i]
					errorMensaje := mw.Partidas[i].Jugadores[j].Conexion.WriteJSON(&enviarMensajeWs)
					if errorMensaje != nil {
						fmt.Println("error al enviar se cierra")
						return
					}
				}
			}
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func (mw *ManejadorWeb) MoverArriba(identificadorPartida, identificadorUsuario string) {
	for i := range mw.Partidas {
		if mw.Partidas[i].Identificador == identificadorPartida {
			for j := range mw.Partidas[i].Jugadores {
				if mw.Partidas[i].Jugadores[j].Identificador == identificadorUsuario {
					if mw.Partidas[i].Jugadores[j].Y > 0 {
						mw.Lock()
						mw.Partidas[i].Jugadores[j].Y = mw.Partidas[i].Jugadores[j].Y - 20
						mw.Unlock()
						var enviarMensajeWs EnviarMensajeWS
						enviarMensajeWs.Mensaje = "jugar"
						enviarMensajeWs.Partida = mw.Partidas[i]

						for k := range mw.Partidas[i].Jugadores {
							mu.Lock()
							mw.Partidas[i].Jugadores[k].Conexion.WriteJSON(&enviarMensajeWs)
							mu.Unlock()
						}
					}
				}
			}
		}
	}
}

func (mw *ManejadorWeb) MoverAbajo(identificadorPartida, identificadorUsuario string) {
	for i := range mw.Partidas {
		if mw.Partidas[i].Identificador == identificadorPartida {
			for j := range mw.Partidas[i].Jugadores {
				if mw.Partidas[i].Jugadores[j].Identificador == identificadorUsuario {
					if mw.Partidas[i].Jugadores[j].Y < 320 {
						mw.Lock()
						mw.Partidas[i].Jugadores[j].Y = mw.Partidas[i].Jugadores[j].Y + 20
						mw.Unlock()
						var enviarMensajeWs EnviarMensajeWS
						enviarMensajeWs.Mensaje = "jugar"
						enviarMensajeWs.Partida = mw.Partidas[i]

						for k := range mw.Partidas[i].Jugadores {
							mu.Lock()
							mw.Partidas[i].Jugadores[k].Conexion.WriteJSON(&enviarMensajeWs)
							mu.Unlock()
						}
					}
				}
			}
		}
	}
}

func (mw *ManejadorWeb) IniciarJuego(identificadorPartida string) {
	for i := range mw.Partidas {
		if mw.Partidas[i].Identificador == identificadorPartida {
			mw.Lock()
			mw.Partidas[i].Jugadores[0].X = 0
			mw.Partidas[i].Jugadores[0].Y = 160
			mw.Partidas[i].Jugadores[0].Ancho = 20
			mw.Partidas[i].Jugadores[0].Alto = 60

			mw.Partidas[i].Jugadores[1].X = 580
			mw.Partidas[i].Jugadores[1].Y = 160
			mw.Partidas[i].Jugadores[1].Ancho = 20
			mw.Partidas[i].Jugadores[1].Alto = 60

			mw.Partidas[i].Pelota.X = 296
			mw.Partidas[i].Pelota.Y = 186

			mw.Unlock()

			var enviarMensajeWs EnviarMensajeWS
			enviarMensajeWs.Mensaje = "jugar"
			enviarMensajeWs.Partida = mw.Partidas[i]
			// enviar escenario
			for j := range mw.Partidas[i].Jugadores {
				conexion := mw.Partidas[i].Jugadores[j].Conexion
				mu.Lock()
				conexion.WriteJSON(&enviarMensajeWs)
				mu.Unlock()
			}

			go mw.MoverPelota(identificadorPartida)
			break
		}
	}
}
