package main

var mutex = NewSemaphore(1)
var oxygen = 0
var hydrogen = 0
var barrier = NewSemaphore(3)
var oxyQueue = NewSemaphore(0)
var hydroQueue = NewSemaphore(0)

type Semaphore struct {
	v    int
	fila chan struct{}
	sc   chan struct{}
}

func NewSemaphore(init int) *Semaphore {
	s := &Semaphore{
		v:    init,
		fila: make(chan struct{}),
		sc:   make(chan struct{}, 1),
	}
	return s
}

func (s *Semaphore) Wait() {
	s.sc <- struct{}{}
	s.v--
	if s.v < 0 {
		<-s.sc
		s.fila <- struct{}{}
	} else {
		<-s.sc
	}
}

func (s *Semaphore) Signal() {
	s.sc <- struct{}{}
	s.v++
	if s.v <= 0 {
		<-s.fila
	}
	<-s.sc
}
