package main

import (
	"fmt"
	"sync"
	"time"
)

var mutex = NewSemaphore(1)
var oxygen = 0
var hydrogen = 0
var barrier = NewBarrier(3)
var oxyQueue = NewSemaphore(0)
var hydroQueue = NewSemaphore(0)

// Estrutura do Semáforo
type Semaphore struct {
	v    int
	fila chan struct{}
	sc   chan struct{}
}

// NewSemaphore inicializa um semáforo
func NewSemaphore(init int) *Semaphore {
	s := &Semaphore{
		v:    init,
		fila: make(chan struct{}),
		sc:   make(chan struct{}, 1),
	}
	return s
}

// Wait diminui o valor do semáforo e bloqueia se necessário
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

// Signal aumenta o valor do semáforo e libera uma thread em espera, se necessário
func (s *Semaphore) Signal() {
	s.sc <- struct{}{}
	s.v++
	if s.v <= 0 {
		<-s.fila
	}
	<-s.sc
}

// Estrutura da Barreira
type Barrier struct {
	n        int
	count    int
	waitCond *sync.Cond
}

// NewBarrier inicializa uma barreira
func NewBarrier(n int) *Barrier {
	return &Barrier{
		n:        n,
		waitCond: sync.NewCond(&sync.Mutex{}),
	}
}

// Wait bloqueia até que a barreira seja atingida por todas as threads
func (b *Barrier) Wait() {
	b.waitCond.L.Lock()
	b.count++
	if b.count == b.n {
		b.count = 0
		b.waitCond.Broadcast()
	} else {
		b.waitCond.Wait()
	}
	b.waitCond.L.Unlock()
}

// Bond simula o processo de ligação para formar uma molécula de água
func Bond() {
	fmt.Println("Bonding...")
	time.Sleep(100 * time.Millisecond) // simula o tempo de ligação
}

// Função da thread Oxigênio
func Oxygen() {
	mutex.Wait()
	oxygen++
	if hydrogen >= 2 {
		hydroQueue.Signal()
		hydroQueue.Signal()
		hydrogen -= 2
		oxyQueue.Signal()
		oxygen--
	} else {
		mutex.Signal()
	}
	oxyQueue.Wait()
	Bond()
	barrier.Wait()
	mutex.Signal()
}

// Função da thread Hidrogênio
func Hydrogen() {
	mutex.Wait()
	hydrogen++
	if hydrogen >= 2 && oxygen >= 1 {
		hydroQueue.Signal()
		hydroQueue.Signal()
		hydrogen -= 2
		oxyQueue.Signal()
		oxygen--
	} else {
		mutex.Signal()
	}
	hydroQueue.Wait()
	Bond()
	barrier.Wait()
}

func main() {
	var wg sync.WaitGroup
	numOxygen := 10
	numHydrogen := 20

	for i := 0; i < numOxygen; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Oxygen()
		}()
	}

	for i := 0; i < numHydrogen; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Hydrogen()
		}()
	}

	wg.Wait()
	fmt.Println("Todas as moléculas de água foram formadas.")
}
