package juego

func (mw *ManejadorWeb) ExistePartida(partidaId string) (bool, int) {
	for i := range mw.Partidas {
		if mw.Partidas[i].Identificador == partidaId {
			return true, len(mw.Partidas[i].Jugadores)
		}
	}
	return false, 0
}

func (mw *ManejadorWeb) AñadirJugadorPartida(jugador Jugador, identificadorSesion string) {
	for i := range mw.Partidas {
		if mw.Partidas[i].Identificador == identificadorSesion {
			mw.Lock()
			mw.Partidas[i].Jugadores = append(mw.Partidas[i].Jugadores, jugador)
			mw.Unlock()
			// AVISAR A LOS 2 JUGADORES
			for j := range mw.Partidas[i].Jugadores {
				var enviarMensajeWs EnviarMensajeWS
				enviarMensajeWs.Mensaje = "segundo jugador"
				mu.Lock()
				mw.Partidas[i].Jugadores[j].Conexion.WriteJSON(&enviarMensajeWs)
				mu.Unlock()
			}
		}
	}
}

func (mw *ManejadorWeb) AñadirPartida(partida Partida) {
	mw.Lock()
	mw.Partidas = append(mw.Partidas, partida)
	mw.Unlock()
}

func (mw *ManejadorWeb) CerrarPartida(identificadorPartida, identificadorUsuario string) {
	for i := range mw.Partidas {
		if mw.Partidas[i].Identificador == identificadorPartida {
			for j := range mw.Partidas[i].Jugadores {
				if mw.Partidas[i].Jugadores[j].Identificador != identificadorUsuario {
					mu.Lock()
					var enviarMensajeWs EnviarMensajeWS
					enviarMensajeWs.Mensaje = "ganador"
					mw.Partidas[i].Jugadores[j].Conexion.WriteJSON(&enviarMensajeWs)
					mu.Unlock()
				}
			}
		}
	}

	mw.Lock()
	for i := range mw.Partidas {
		if mw.Partidas[i].Identificador == identificadorPartida {
			mw.Partidas = append(mw.Partidas[:i], mw.Partidas[i+1:]...)
			break
		}
	}
	mw.Unlock()
}
