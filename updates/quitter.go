package updates

import "sync"

type Quitter struct {
	lock   *sync.Mutex
	quit   chan struct{}
	active bool
}

func NewQuitter() (quitter *Quitter) {
	quitter = &Quitter{lock: &sync.Mutex{}}
	quitter.Activate()
	return
}

func (q *Quitter) Activate() {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.active {
		return
	}
	q.quit = make(chan struct{})
	q.active = true
}

func (q *Quitter) Active() bool {
	return q.active
}

func (q *Quitter) Quitter() <-chan struct{} {
	return q.quit
}

func (q *Quitter) Quit() {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.active {
		close(q.quit)
		q.active = false
	}
}
