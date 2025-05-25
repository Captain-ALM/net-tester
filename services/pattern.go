package services

type PatternService struct {
	quit          chan struct{}
	Pattern       string
	MinLength     uint
	MaxLength     uint
	currentLength uint
}

func (p *PatternService) init(quit chan struct{}) {
	p.quit = quit
	p.currentLength = p.MinLength
}

func (p *PatternService) write(bytes []byte) {
	//TODO implement me
	panic("implement me")
}

func (p *PatternService) read() []byte {
	//TODO implement me
	panic("implement me")
}
