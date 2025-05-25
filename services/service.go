package services

import "gitcove.com/alfred/net-tester/updates"

type Service interface {
	Init(quit *updates.Quitter, updater *updates.Update)
	Write([]byte)
	Read() []byte
}
