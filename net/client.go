package net

import (
	"gitcove.com/alfred/net-tester/services"
	"gitcove.com/alfred/net-tester/updates"
	"net"
	"time"
)

func RunClient(conn net.Conn, service services.Service, quitter *updates.Quitter, updater *updates.Update, bufferSize uint, timeout time.Duration) {
	service.Init(quitter, updater)
	go func() {
		select {
		case <-quitter.Quitter():
			_ = conn.Close()
		}
	}()
	go func() {
		bytes := make([]byte, bufferSize)
		for quitter.Active() {
			if timeout > 0 {
				_ = conn.SetReadDeadline(time.Now().Add(timeout))
			}
			n, err := conn.Read(bytes)
			updater.BytesReceived += uint64(n)
			if err != nil {
				quitter.Quit()
				return
			}
			service.Write(bytes[:n])
		}
	}()
	go func() {
		for quitter.Active() {
			if timeout > 0 {
				_ = conn.SetWriteDeadline(time.Now().Add(timeout))
			}
			n, err := conn.Write(service.Read())
			updater.BytesSent += uint64(n)
			if err != nil {
				quitter.Quit()
				return
			}
		}
	}()
}
