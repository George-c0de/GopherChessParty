package utils

import "github.com/gorilla/websocket"

func CloseWebsocket(Conn *websocket.Conn) error {
	err := Conn.Close()
	if err != nil {
		return err
	}
	return nil
}
