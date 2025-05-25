package services

type Service interface {
	init(quit chan struct{})
	write([]byte)
	read() []byte
}
