package services

import (
	bytes2 "bytes"
	"gitcove.com/alfred/net-tester/updates"
)

type PatternService struct {
	quit               *updates.Quitter
	Pattern            []byte
	cached             []byte
	MinLength          uint
	MaxLength          uint
	currentLengthWrite uint
	currentLengthRead  uint
	update             *updates.Update
}

func (p *PatternService) Init(quit *updates.Quitter, updater *updates.Update) {
	p.quit = quit
	p.currentLengthWrite = p.MinLength
	p.currentLengthRead = p.MinLength
	p.update = updater
}

func (p *PatternService) Write(bytes []byte) {
	if p.currentLengthWrite > p.MaxLength {
		p.quit.Quit()
		return
	}
	if p.quit.Active() {
		rBytes := make([]byte, 0, len(bytes)+len(p.cached))
		rBytes = append(rBytes, p.cached...)
		rBytes = append(rBytes, bytes...)
		for len(rBytes)-int(p.currentLengthWrite) >= 0 && p.currentLengthWrite <= p.MaxLength {
			cBytes := rBytes[:p.currentLengthWrite]
			rBytes = rBytes[p.currentLengthWrite:]
			p.update.PatternLengthIn = p.currentLengthWrite
			if !bytes2.Equal(cBytes, p.getDataFromPattern(p.currentLengthWrite)) {
				p.quit.Quit()
				return
			}
			p.currentLengthWrite++
		}
		p.cached = rBytes
	}
	if p.currentLengthWrite >= p.MaxLength {
		p.quit.Quit()
	}
}

func (p *PatternService) Read() (read []byte) {
	if p.quit.Active() {
		if p.currentLengthRead > p.MaxLength {
			read = nil
			return
		}
		read = p.getDataFromPattern(p.currentLengthRead)
		p.update.PatternLengthOut = p.currentLengthRead
		p.currentLengthRead++
	}
	return
}

func (p *PatternService) getDataFromPattern(l uint) (toReturn []byte) {
	toReturn = make([]byte, 0, l)
	for l > 0 {
		cl := min(uint(len(p.Pattern)), l)
		l -= cl
		toReturn = append(toReturn, p.Pattern[:cl]...)
	}
	return
}
