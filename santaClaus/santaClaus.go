package main

import (
	"fmt"
	"time"
)

var elves = 0
var reindeer = 0
var mutex = NewSemaphore(1)
var santaSem = NewSemaphore(0)
var reindeerSem = NewSemaphore(0)
var elfTex = NewSemaphore(1)

func main() {

	for i := 1; i <= 1000; i++ {
		go elfArrives()
	}

	for i := 1; i <= 9; i++ {
		go reindeerArrives()
	}

	for {
		santaSem.Wait()
		mutex.Wait()

		if reindeer == 9 {
			prepareSleigh()
			reindeerSem.Signal()
			fmt.Println("É Natal, entregando presentes!")
			return
		} else if elves == 3 {
			helpElves()
		}

		mutex.Signal()
	}
}

func elfArrives() {
	elfTex.Wait()
	mutex.Wait()
	elves++

	fmt.Println("Elfo: Pede Ajuda", elves)
	if elves == 3 {
		santaSem.Signal()

	} else {

		elfTex.Signal()

	}
	mutex.Signal()

	getHelp()

	mutex.Wait()
	elves--
	if elves == 0 {
		elfTex.Signal()
	}
	mutex.Signal()
}

func reindeerArrives() {
	mutex.Wait()
	reindeer++
	fmt.Println("Rena: Volta das férias", reindeer)

	if reindeer == 9 {
		for i := 0; i < 9; i++ {
			santaSem.Signal()
		}
	}
	mutex.Signal()

	reindeerSem.Wait()
	getHitched()
}

func prepareSleigh() {
	fmt.Println("Santa: Preparando o trenó")
}

func getHitched() {
	fmt.Println("Santa: Preparando as renas")
	reindeerSem.Signal()
}

func helpElves() {
	fmt.Println("Santa: Ajudando os elfos", elves)
	time.Sleep(1 * time.Second)
}

func getHelp() {
	fmt.Println("Elfos: Esperando ajuda", elves)
	time.Sleep(1 * time.Second)
}

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
