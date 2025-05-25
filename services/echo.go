package services

import "gitcove.com/alfred/net-tester/updates"

type EchoService struct {
	forward chan []byte
	quit    *updates.Quitter
}

func (e *EchoService) Init(quit *updates.Quitter, _ *updates.Update) {
	e.quit = quit
	e.forward = make(chan []byte)
}

func (e *EchoService) Write(bytes []byte) {
	cBytes := make([]byte, len(bytes))
	copy(cBytes, bytes)
	select {
	case e.forward <- cBytes:
	case <-e.quit.Quitter():
	}
}

func (e *EchoService) Read() (bytes []byte) {
	select {
	case bytes = <-e.forward:
	case <-e.quit.Quitter():
		bytes = nil
	}
	return
}
