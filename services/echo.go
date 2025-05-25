package services

type EchoService struct {
	forward chan []byte
	quit    chan struct{}
}

func (e *EchoService) init(quit chan struct{}) {
	e.quit = quit
	e.forward = make(chan []byte)
}

func (e *EchoService) write(bytes []byte) {
	select {
	case e.forward <- bytes:
	case <-e.quit:
	}
}

func (e *EchoService) read() (bytes []byte) {
	select {
	case bytes = <-e.forward:
	case <-e.quit:
		bytes = make([]byte, 0)
	}
	return
}
